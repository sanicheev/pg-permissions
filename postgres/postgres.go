package postgres

import (
	"encoding/json"
	"html/template"
	"os"
	"pg_permissions/database"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
)

const (
	REQUEST_MULIPLIER = 4
)

var (
	templ *template.Template
)

type Controller struct {
	dbController   *database.Controller
	permissionType string
	maxRequests    int
	reportHTML     bool
	reportJSON     bool

	mutex *sync.Mutex

	users       []database.User
	databases   []database.Database
	schemas     []database.Schema
	tables      []database.Table
	columns     []database.Column
	functions   []database.Function
	languages   []database.Language
	types       []database.Type
	sequences   []database.Sequence
	fdws        []database.FDW
	servers     []database.Server
	tablespaces []database.Tablespace
	roles       []database.Role

	userPermissionMap UserPermissionMap
}

func (c *Controller) Lock() {
	c.mutex.Lock()
}

func (c *Controller) Unlock() {
	c.mutex.Unlock()
}

func NewController(dbc *database.Controller, permissionType string, maxRequests int, reportHTML bool, reportJSON bool) *Controller {
	postgresController := Controller{
		dbController:   dbc,
		permissionType: permissionType,
		mutex:          new(sync.Mutex),
		maxRequests:    maxRequests,
		reportHTML:     reportHTML,
		reportJSON:     reportJSON,
	}
	return &postgresController
}

func (c *Controller) Run() error {
	maxRequests := make(chan int, REQUEST_MULIPLIER)

	gatheringObjectsFuncs := []func() error{
		c.listUsers,
		c.listColumns,
		c.listDatabases,
		c.listSchemas,
		c.listTables,
		c.listFunctions,
		c.listLanguages,
		c.listTypes,
		c.listFDWs,
		c.listServers,
		c.listTablespaces,
		c.listRoles,
	}

	for idx, _ := range gatheringObjectsFuncs {
		err := gatheringObjectsFuncs[idx]()
		if err != nil {
			return err
		}
	}

	userPermissionMap := make(UserPermissionMap)

	gatheringPermissionsFuncs := []func(UserPermissionMap, *string) error{
		c.getDatabasePermissions,
		c.getSchemaPermissions,
		c.getTablePermissions,
		c.getColumnPermissions,
		c.getFunctionPermissions,
		c.getLanguagePermissions,
		c.getTypePermissions,
		c.getSequencePermissions,
		c.getFDWPermissions,
		c.getServerPermissions,
		c.getTablespacePermissions,
		c.getRolePermissions,
	}

	g := new(errgroup.Group)

	for idx_u, _ := range c.users {
		maxRequests <- 1
		idx_u := idx_u

		g.Go(func() error {
			objectMap := ObjectMap{
				FailedToCheck: make(map[string]string),
				ObjectDesc:    make(map[string]ObjectDesc),
				Permissions:   make(map[string][]string),
			}

			c.Lock()
			userPermissionMap[c.users[idx_u].Name] = objectMap
			c.Unlock()

			for idx_f, _ := range gatheringPermissionsFuncs {
				err := gatheringPermissionsFuncs[idx_f](userPermissionMap, &c.users[idx_u].Name)
				if err != nil {
					// should we just log the error
					<-maxRequests
					return err
				}
			}

			<-maxRequests
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	if c.reportHTML {
		funcMap := template.FuncMap{
			"Title": strings.Title,
		}

		tmpl, err := template.New("results.html.tpl").Funcs(funcMap).ParseFiles("results.html.tpl")
		if err != nil {
			return err
		}

		fh, err := os.OpenFile("results.html", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		defer fh.Close()
		if err != nil {
			return err
		}

		err = tmpl.Execute(fh, userPermissionMap)
		if err != nil {
			return err
		}
	}

	if c.reportJSON {
		jsonReport, err := json.Marshal(userPermissionMap)
		if err != nil {
			return err
		}

		fj, err := os.OpenFile("results.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		defer fj.Close()
		if err != nil {
			return err
		}

		_, err = fj.Write(jsonReport)
		if err != nil {
			return err
		}
	}

	return nil
}
