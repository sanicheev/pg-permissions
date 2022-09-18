package database

func (c *Controller) ListDatabases() ([]Database, error) {
	databases := []Database{}

	err := c.Database.Select(
		&databases,
		`SELECT
			datname AS name,
			datistemplate AS is_template,
			datallowconn AS allow_connect
                FROM pg_catalog.pg_database`,
	)

	return databases, err
}

func (c *Controller) FetchDatabasePrivilege(user *string, database *string) (*PrivilegeDescriptor, error) {
	privilegeDescriptor := PrivilegeDescriptor{}

	err := c.Database.Get(
		&privilegeDescriptor,
		`SELECT
			has_database_privilege($1, $2, 'CONNECT') AS connect,
			has_database_privilege($1, $2, 'CREATE') AS create,
			has_database_privilege($1, $2, 'TEMPORARY') AS temporary;`,
		*user,
		*database,
	)

	if err != nil {
		return nil, err
	}

	return &privilegeDescriptor, nil
}

func (c *Controller) DatabasePermissions(permissionType *string) []string {
	R := []string{
		"CONNECT",
	}
	W := []string{
		"CREATE",
		"TEMPORARY",
	}

	ALL := append(W, R...)

	if *permissionType == "W" {
		return W
	}

	return ALL
}
