package postgres

import (
	"golang.org/x/sync/errgroup"

	"pg_permissions/log"
	"time"
)

func (c *Controller) listTablespaces() error {
	log.Debugln("Fetching tablespace list from database")

	tablespaces, err := c.dbController.ListTablespaces()

	if err != nil {
		return err
	}

	c.tablespaces = tablespaces

	return nil

}

func (c *Controller) getTablespacePermissions(upm UserPermissionMap, username *string) error {
	if c.permissionType == "R" {
		return nil
	}

	maxRequests := make(chan int, c.maxRequests)

	objectDesc := make(ObjectDesc)

	permissions := c.dbController.TablespacePermissions(&c.permissionType)

	c.Lock()
	upm[*username].ObjectDesc["tablespace"] = objectDesc
	upm[*username].Permissions["tablespace"] = permissions
	c.Unlock()

	g := new(errgroup.Group)

	for idx, _ := range c.tablespaces {
		maxRequests <- 1
		idx := idx

		g.Go(func() error {
			time.Sleep(1 * time.Second)

			privilegeDescriptor, err := c.dbController.FetchTablespacePrivilege(username, &c.tablespaces[idx].Name)
			if err != nil {
				<-maxRequests
				return err
			}
			permissionMap := PrivilegeDescriptorToPermissionMap(privilegeDescriptor, permissions)

			c.Lock()
			upm[*username].ObjectDesc["tablespace"][c.tablespaces[idx].Name] = permissionMap
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
