package postgres

import (
	"fmt"
	"pg_permissions/log"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/lib/pq"
)

func (c *Controller) listFunctions() error {
	log.Debugln("Fetching function list from database")

	functions, err := c.dbController.ListFunctions()

	if err != nil {
		return err
	}

	c.functions = functions

	return nil

}

func (c *Controller) getFunctionPermissions(upm UserPermissionMap, username *string) error {
	if c.permissionType == "W" {
		return nil
	}

	maxRequests := make(chan int, c.maxRequests)

	objectDesc := make(ObjectDesc)

	permissions := c.dbController.FunctionPermissions(&c.permissionType)

	c.Lock()
	upm[*username].ObjectDesc["function"] = objectDesc
	upm[*username].Permissions["function"] = permissions
	c.Unlock()

	g := new(errgroup.Group)

	for idx, _ := range c.functions {
		maxRequests <- 1
		idx := idx

		g.Go(func() error {
			time.Sleep(1 * time.Second)

			c.functions[idx].SanitizeArguments()

			functionNameWithArgs := fmt.Sprintf(
				"%s.%s(%s)",
				c.functions[idx].Schema,
				c.functions[idx].Name,
				c.functions[idx].Arguments,
			)

			c.Lock()
			_, processed := upm[*username].ObjectDesc["function"][functionNameWithArgs]
			_, skipped := upm[*username].FailedToCheck[functionNameWithArgs]
			c.Unlock()
			if processed || skipped {
				log.Debugf("Function: %v has already been processed!", functionNameWithArgs)
				<-maxRequests
				return nil
			}

			privilegeDescriptor, err := c.dbController.FetchFunctionPrivilege(username, &functionNameWithArgs)
			if err != nil {
				if err.(*pq.Error).Code.Name() == "undefined_function" {
					log.Warnf("Could not find a function: %v. Permission issue? Please check it manually.", functionNameWithArgs)
					c.Lock()
					upm[*username].FailedToCheck[functionNameWithArgs] = "function"
					c.Unlock()
					<-maxRequests
					return nil
				} else {
					<-maxRequests
					return err
				}
			}

			permissionMap := PrivilegeDescriptorToPermissionMap(privilegeDescriptor, permissions)

			c.Lock()
			upm[*username].ObjectDesc["function"][functionNameWithArgs] = permissionMap
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
