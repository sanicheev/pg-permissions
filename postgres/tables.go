package postgres

import (
	"golang.org/x/sync/errgroup"

	"fmt"
	"pg_permissions/log"
	"time"
)

func (c *Controller) listTables() error {
	log.Debugln("Fetching table list from database")

	tables, err := c.dbController.ListTables()

	if err != nil {
		return err
	}

	c.tables = tables

	return nil

}

func (c *Controller) getTablePermissions(upm UserPermissionMap, username *string) error {
	maxRequests := make(chan int, c.maxRequests)

	objectDesc := make(ObjectDesc)

	permissions := c.dbController.TablePermissions(&c.permissionType)

	c.Lock()
	upm[*username].ObjectDesc["table"] = objectDesc
	upm[*username].Permissions["table"] = permissions
	c.Unlock()

	g := new(errgroup.Group)

	for idx, _ := range c.tables {
		maxRequests <- 1
		idx := idx

		g.Go(func() error {
			time.Sleep(1 * time.Second)

			fullTableName := fmt.Sprintf("%s.%s", c.tables[idx].Schema, c.tables[idx].Name)

			privilegeDescriptor, err := c.dbController.FetchTablePrivilege(username, &fullTableName)
			if err != nil {
				<-maxRequests
				return err
			}
			permissionMap := PrivilegeDescriptorToPermissionMap(privilegeDescriptor, permissions)

			c.Lock()
			upm[*username].ObjectDesc["table"][fullTableName] = permissionMap
			c.Unlock()

			<-maxRequests
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
