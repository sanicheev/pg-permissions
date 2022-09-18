package postgres

import (
	"pg_permissions/log"
)

func (c *Controller) listUsers() error {
	log.Debugln("Fetching user list from database")

	users, err := c.dbController.ListUsers()

	if err != nil {
		return err
	}

	c.users = users

	return nil
}
