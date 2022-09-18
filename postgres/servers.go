package postgres

import (
	"golang.org/x/sync/errgroup"

	"pg_permissions/log"
	"time"
)

func (c *Controller) listServers() error {
	log.Debugln("Fetching server list from database")

	servers, err := c.dbController.ListServers()
	c.servers = servers

	return err

}

func (c *Controller) getServerPermissions(upm UserPermissionMap, username *string) error {
	if c.permissionType == "W" {
		return nil
	}

	maxRequests := make(chan int, c.maxRequests)

	objectDesc := make(ObjectDesc)

	permissions := c.dbController.ServerPermissions(&c.permissionType)

	c.Lock()
	upm[*username].ObjectDesc["server"] = objectDesc
	upm[*username].Permissions["server"] = permissions
	c.Unlock()

	g := new(errgroup.Group)

	for idx, _ := range c.servers {
		maxRequests <- 1
		idx := idx

		g.Go(func() error {
			time.Sleep(1 * time.Second)

			privilegeDescriptor, err := c.dbController.FetchServerPrivilege(username, &c.servers[idx].Name)
			if err != nil {
				<-maxRequests
				return err
			}
			permissionMap := PrivilegeDescriptorToPermissionMap(privilegeDescriptor, permissions)

			c.Lock()
			upm[*username].ObjectDesc["server"][c.servers[idx].Name] = permissionMap
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
