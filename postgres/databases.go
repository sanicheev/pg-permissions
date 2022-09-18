package postgres

import (
	"pg_permissions/log"

	"golang.org/x/sync/errgroup"

	"time"
)

func (c *Controller) listDatabases() error {
	log.Debugln("Fetching database list from database")

	databases, err := c.dbController.ListDatabases()

	if err != nil {
		return err
	}

	c.databases = databases

	return nil

}

func (c *Controller) getDatabasePermissions(upm UserPermissionMap, username *string) error {
	maxRequests := make(chan int, c.maxRequests)

	objectDesc := make(ObjectDesc)

	permissions := c.dbController.DatabasePermissions(&c.permissionType)

	c.Lock()
	upm[*username].ObjectDesc["database"] = objectDesc
	upm[*username].Permissions["database"] = permissions
	c.Unlock()

	g := new(errgroup.Group)

	for idx, _ := range c.databases {
		maxRequests <- 1
		idx := idx

		g.Go(func() error {
			time.Sleep(1 * time.Second)

			privilegeDescriptor, err := c.dbController.FetchDatabasePrivilege(username, &c.databases[idx].Name)
			if err != nil {
				<-maxRequests
				return err
			}
			permissionMap := PrivilegeDescriptorToPermissionMap(privilegeDescriptor, permissions)

			c.Lock()
			upm[*username].ObjectDesc["database"][c.databases[idx].Name] = permissionMap
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
