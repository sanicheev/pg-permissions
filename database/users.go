package database

func (c *Controller) ListUsers() ([]User, error) {
	users := []User{}

	err := c.Database.Select(
		&users,
		`SELECT
			usename AS name,
			usecreatedb AS createdb_privilege,
			usesuper AS superuser_privilege,
			userepl AS replication_privilege,
			usebypassrls AS bypass_row_level_security_policy,
			valuntil AS password_valid_until
		FROM pg_catalog.pg_user`,
	)

	return users, err
}
