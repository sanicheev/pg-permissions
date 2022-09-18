package database

func (c *Controller) ListTablespaces() ([]Tablespace, error) {
	tablespaces := []Tablespace{}

	err := c.Database.Select(
		&tablespaces,
		`SELECT
			spcname AS name
		FROM pg_catalog.pg_tablespace`,
	)

	return tablespaces, err
}

func (c *Controller) FetchTablespacePrivilege(user *string, tablespace *string) (*PrivilegeDescriptor, error) {
	privilegeDescriptor := PrivilegeDescriptor{}

	err := c.Database.Get(
		&privilegeDescriptor,
		`SELECT
			has_tablespace_privilege($1, $2, 'CREATE') AS create;`,
		*user,
		*tablespace,
	)

	if err != nil {
		return nil, err
	}

	return &privilegeDescriptor, nil
}

func (c *Controller) TablespacePermissions(permissionType *string) []string {
	R := []string{}

	W := []string{
		"CREATE",
	}

	ALL := append(W, R...)

	if *permissionType == "W" {
		return W
	}

	return ALL
}
