package postgres

import (
	"golang.org/x/sync/errgroup"

	"pg_permissions/log"

	"time"
)

func (c *Controller) listFDWs() error {
	log.Debugln("Fetching fdw list from database")

	fdws, err := c.dbController.ListFDWs()

	if err != nil {
		return err
	}

	c.fdws = fdws

	return nil

}
func (c *Controller) getFDWPermissions(upm UserPermissionMap, username *string) error {
	if c.permissionType == "W" {
		return nil
	}

	maxRequests := make(chan int, c.maxRequests)

	objectDesc := make(ObjectDesc)

	permissions := c.dbController.FDWPermissions(&c.permissionType)

	c.Lock()
	upm[*username].ObjectDesc["fdw"] = objectDesc
	upm[*username].Permissions["fdw"] = permissions
	c.Unlock()

	g := new(errgroup.Group)

	for idx, _ := range c.fdws {
		maxRequests <- 1
		idx := idx

		g.Go(func() error {
			time.Sleep(1 * time.Second)

			privilegeDescriptor, err := c.dbController.FetchFDWPrivilege(username, &c.fdws[idx].Name)
			if err != nil {
				<-maxRequests
				return err
			}
			permissionMap := PrivilegeDescriptorToPermissionMap(privilegeDescriptor, permissions)

			c.Lock()
			upm[*username].ObjectDesc["fdw"][c.fdws[idx].Name] = permissionMap
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
