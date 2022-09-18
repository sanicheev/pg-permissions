package postgres

import (
	"golang.org/x/sync/errgroup"

	"pg_permissions/log"
	"time"
)

func (c *Controller) listSchemas() error {
	log.Debugln("Fetching schema list from database")

	schemas, err := c.dbController.ListSchemas()

	if err != nil {
		return err
	}

	c.schemas = schemas

	return nil

}

func (c *Controller) getSchemaPermissions(upm UserPermissionMap, username *string) error {
	maxRequests := make(chan int, c.maxRequests)

	objectDesc := make(ObjectDesc)

	permissions := c.dbController.SchemaPermissions(&c.permissionType)

	c.Lock()
	upm[*username].ObjectDesc["schema"] = objectDesc
	upm[*username].Permissions["schema"] = permissions
	c.Unlock()

	g := new(errgroup.Group)

	for idx, _ := range c.schemas {
		maxRequests <- 1
		idx := idx

		g.Go(func() error {
			time.Sleep(1 * time.Second)

			privilegeDescriptor, err := c.dbController.FetchSchemaPrivilege(username, &c.schemas[idx].Name)
			if err != nil {
				<-maxRequests
				return err
			}
			permissionMap := PrivilegeDescriptorToPermissionMap(privilegeDescriptor, permissions)

			c.Lock()
			upm[*username].ObjectDesc["schema"][c.schemas[idx].Name] = permissionMap
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
