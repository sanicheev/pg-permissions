package postgres

import (
	"fmt"
	"pg_permissions/log"

	"golang.org/x/sync/errgroup"

	"time"
)

func (c *Controller) listSequences() error {
	log.Debugln("Fetching sequence list from database")

	sequences, err := c.dbController.ListSequences()

	if err != nil {
		return err
	}

	c.sequences = sequences

	return nil

}

func (c *Controller) getSequencePermissions(upm UserPermissionMap, username *string) error {
	maxRequests := make(chan int, c.maxRequests)

	objectDesc := make(ObjectDesc)

	permissions := c.dbController.SequencePermissions(&c.permissionType)

	c.Lock()
	upm[*username].ObjectDesc["sequence"] = objectDesc
	upm[*username].Permissions["sequence"] = permissions
	c.Unlock()

	g := new(errgroup.Group)

	for idx, _ := range c.sequences {
		maxRequests <- 1
		idx := idx

		g.Go(func() error {
			time.Sleep(1 * time.Second)

			fullSequenceName := fmt.Sprintf("%s.%s", c.sequences[idx].Schema, c.sequences[idx].Name)

			privilegeDescriptor, err := c.dbController.FetchSequencePrivilege(username, &fullSequenceName)
			if err != nil {
				<-maxRequests
				return err
			}
			permissionMap := PrivilegeDescriptorToPermissionMap(privilegeDescriptor, permissions)

			c.Lock()
			upm[*username].ObjectDesc["sequence"][fullSequenceName] = permissionMap
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
