package postgres

import (
	"fmt"
	"pg_permissions/log"

	"golang.org/x/sync/errgroup"

	"time"
)

func (c *Controller) listColumns() error {
	log.Debugln("Fetching column list from database")

	columns, err := c.dbController.ListColumns()

	if err != nil {
		return err
	}

	c.columns = columns

	return nil
}

func (c *Controller) getColumnPermissions(upm UserPermissionMap, username *string) error {
	maxRequests := make(chan int, c.maxRequests)

	objectDesc := make(ObjectDesc)
	permissions := c.dbController.ColumnPermissions(&c.permissionType)

	c.Lock()
	upm[*username].ObjectDesc["column"] = objectDesc
	upm[*username].Permissions["column"] = permissions
	c.Unlock()

	g := new(errgroup.Group)

	for idx, _ := range c.columns {
		maxRequests <- 1
		idx := idx

		g.Go(func() error {
			time.Sleep(1 * time.Second)

			fullTableName := fmt.Sprintf("%s.%s", c.columns[idx].Schema, c.columns[idx].Table)
			columnName := fmt.Sprintf("%s:%s", fullTableName, c.columns[idx].Name)

			privilegeDescriptor, err := c.dbController.FetchColumnPrivilege(username, &fullTableName, &c.columns[idx].Name)
			if err != nil {
				<-maxRequests
				return err
			}
			permissionMap := PrivilegeDescriptorToPermissionMap(privilegeDescriptor, permissions)

			c.Lock()
			upm[*username].ObjectDesc["column"][columnName] = permissionMap
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
