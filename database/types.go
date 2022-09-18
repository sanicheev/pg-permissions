package database

func (c *Controller) ListTypes() ([]Type, error) {
	types := []Type{}

	err := c.Database.Select(
		&types,
		`SELECT
			pgtype.typname AS name,
			pgnamespace.nspname AS schema
		FROM pg_catalog.pg_type AS pgtype
		INNER JOIN pg_catalog.pg_namespace pgnamespace ON pgnamespace.oid = pgtype.typnamespace
		WHERE pgnamespace.nspname != 'pg_toast';`,
	)

	return types, err
}

func (c *Controller) FetchTypePrivilege(user *string, typeName *string) (*PrivilegeDescriptor, error) {
	privilegeDescriptor := PrivilegeDescriptor{}

	err := c.Database.Get(
		&privilegeDescriptor,
		`SELECT
			has_type_privilege($1, $2, 'USAGE') AS usage;`,
		*user,
		*typeName,
	)

	if err != nil {
		return nil, err
	}

	return &privilegeDescriptor, nil
}

func (c *Controller) TypePermissions(permissionType *string) []string {
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
