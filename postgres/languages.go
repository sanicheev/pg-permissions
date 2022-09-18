package postgres

import (
	"golang.org/x/sync/errgroup"

	"pg_permissions/log"
	"time"
)

func (c *Controller) listLanguages() error {
	log.Debugln("Fetching language list from database")

	languages, err := c.dbController.ListLanguages()

	if err != nil {
		return err
	}

	c.languages = languages

	return nil

}

func (c *Controller) getLanguagePermissions(upm UserPermissionMap, username *string) error {
	if c.permissionType == "W" {
		return nil
	}

	maxRequests := make(chan int, c.maxRequests)

	objectDesc := make(ObjectDesc)

	permissions := c.dbController.LanguagePermissions(&c.permissionType)

	c.Lock()
	upm[*username].ObjectDesc["language"] = objectDesc
	upm[*username].Permissions["language"] = permissions
	c.Unlock()

	g := new(errgroup.Group)

	for idx, _ := range c.languages {
		maxRequests <- 1
		idx := idx

		g.Go(func() error {
			time.Sleep(1 * time.Second)

			privilegeDescriptor, err := c.dbController.FetchLanguagePrivilege(username, &c.languages[idx].Language)
			if err != nil {
				<-maxRequests
				return err
			}
			permissionMap := PrivilegeDescriptorToPermissionMap(privilegeDescriptor, permissions)

			c.Lock()
			upm[*username].ObjectDesc["language"][c.languages[idx].Language] = permissionMap
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
