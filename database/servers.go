package database

func (c *Controller) ListServers() ([]Server, error) {
	servers := []Server{}

	err := c.Database.Select(
		&servers,
		`SELECT
			srvname AS name,
			srvfdw AS fdw
                FROM pg_catalog.pg_foreign_server`,
	)

	return servers, err
}

func (c *Controller) FetchServerPrivilege(user *string, server *string) (*PrivilegeDescriptor, error) {
	privilegeDescriptor := PrivilegeDescriptor{}

	err := c.Database.Get(
		&privilegeDescriptor,
		`SELECT
			has_server_privilege($1, $2, 'USAGE') AS usage;`,
		*user,
		*server,
	)

	if err != nil {
		return nil, err
	}

	return &privilegeDescriptor, nil
}

func (c *Controller) ServerPermissions(permissionType *string) []string {
	R := []string{
		"USAGE",
	}

	W := []string{}

	ALL := append(W, R...)

	if *permissionType == "W" {
		return W
	}

	return ALL
}
