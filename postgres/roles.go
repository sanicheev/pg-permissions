package postgres

import (
	"golang.org/x/sync/errgroup"

	"pg_permissions/log"
	"time"
)

func (c *Controller) listRoles() error {
	log.Debugln("Fetching role list from database")

	roles, err := c.dbController.ListRoles()

	if err != nil {
		return err
	}

	c.roles = roles

	return nil

}

func (c *Controller) getRolePermissions(upm UserPermissionMap, username *string) error {
	maxRequests := make(chan int, c.maxRequests)

	objectDesc := make(ObjectDesc)

	usagePermissions := c.dbController.RoleUsagePermissions()

	rolePermissions := c.dbController.RolePermissions(&c.permissionType)

	c.Lock()
	upm[*username].ObjectDesc["role"] = objectDesc
	upm[*username].Permissions["role"] = rolePermissions
	c.Unlock()

	g := new(errgroup.Group)

	for idx, _ := range c.roles {
		maxRequests <- 1
		idx := idx

		g.Go(func() error {
			time.Sleep(1 * time.Second)

			canUseRole := false
			privilegeDescriptor, err := c.dbController.FetchRolePrivilege(username, &c.roles[idx].Name)
			if err != nil {
				<-maxRequests
				return err
			}
			usagePermissionMap := PrivilegeDescriptorToPermissionMap(privilegeDescriptor, usagePermissions)

			for _, v := range usagePermissionMap {
				if v == true {
					canUseRole = true
					break
				}
			}

			if canUseRole == true {
				permissionMap := PrivilegeDescriptorToRolePermissionMap(&c.roles[idx], rolePermissions)
				c.Lock()
				upm[*username].ObjectDesc["role"][c.roles[idx].Name] = permissionMap
				c.Unlock()
			}

			<-maxRequests
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
