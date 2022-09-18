package database

func (c *Controller) ListFDWs() ([]FDW, error) {
	fdws := []FDW{}

	err := c.Database.Select(
		&fdws,
		`SELECT
			fdwname AS name
                FROM pg_catalog.pg_foreign_data_wrapper`,
	)

	return fdws, err
}

func (c *Controller) FetchFDWPrivilege(user *string, fdw *string) (*PrivilegeDescriptor, error) {
	privilegeDescriptor := PrivilegeDescriptor{}

	err := c.Database.Get(
		&privilegeDescriptor,
		`SELECT
			has_foreign_data_wrapper_privilege($1, $2, 'USAGE') AS usage;`,
		*user,
		*fdw,
	)

	if err != nil {
		return nil, err
	}

	return &privilegeDescriptor, nil
}

func (c *Controller) FDWPermissions(permissionType *string) []string {
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
