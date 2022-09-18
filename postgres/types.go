package postgres

import (
	"golang.org/x/sync/errgroup"

	"fmt"
	"pg_permissions/log"
	"time"
)

func (c *Controller) listTypes() error {
	log.Debugln("Fetching types list from database")

	types, err := c.dbController.ListTypes()

	if err != nil {
		return err
	}

	c.types = types

	return err

}

func (c *Controller) getTypePermissions(upm UserPermissionMap, username *string) error {
	if c.permissionType == "W" {
		return nil
	}

	maxRequests := make(chan int, c.maxRequests)

	objectDesc := make(ObjectDesc)

	permissions := c.dbController.TypePermissions(&c.permissionType)

	c.Lock()
	upm[*username].ObjectDesc["type"] = objectDesc
	upm[*username].Permissions["type"] = permissions
	c.Unlock()

	g := new(errgroup.Group)

	for idx, _ := range c.types {
		maxRequests <- 1
		idx := idx

		g.Go(func() error {
			time.Sleep(1 * time.Second)

			fullTypeName := fmt.Sprintf("%s.%s", c.types[idx].Schema, c.types[idx].Name)

			privilegeDescriptor, err := c.dbController.FetchTypePrivilege(username, &fullTypeName)
			if err != nil {
				<-maxRequests
				return err
			}
			permissionMap := PrivilegeDescriptorToPermissionMap(privilegeDescriptor, permissions)

			c.Lock()
			upm[*username].ObjectDesc["type"][fullTypeName] = permissionMap
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
