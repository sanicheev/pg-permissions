package database

func (c *Controller) ListRoles() ([]Role, error) {
	roles := []Role{}

	err := c.Database.Select(
		&roles,
		`SELECT
			rolname AS name,
			rolsuper AS superuser,
			rolcreaterole AS createrole,
			rolcreatedb AS createdb,
			rolcanlogin AS canlogin,
			rolreplication AS replication
                FROM pg_catalog.pg_roles`,
	)

	return roles, err
}

func (c *Controller) FetchRolePrivilege(user *string, role *string) (*PrivilegeDescriptor, error) {
	privilegeDescriptor := PrivilegeDescriptor{}

	err := c.Database.Get(
		&privilegeDescriptor,
		`SELECT
			pg_has_role($1, $2, 'USAGE') AS usage,
			pg_has_role($1, $2, 'MEMBER') AS member;`,
		*user,
		*role,
	)

	if err != nil {
		return nil, err
	}

	return &privilegeDescriptor, nil
}

func (c *Controller) RoleUsagePermissions() []string {
	return []string{
		"USAGE",  // whether the privileges of the role are immediately available (without doing SET ROLE).
		"MEMBER", // whether the user has direct or indirect membership in the role (the right to do SET ROLE).
	}
}

func (c *Controller) RolePermissions(permissionType *string) []string {
	R := []string{
		"CANLOGIN",
		"REPLICATION",
	}

	W := []string{
		"SUPERUSER",
		"CREATEROLE",
		"CREATEDB",
	}

	ALL := append(W, R...)

	if *permissionType == "W" {
		return W
	}

	return ALL
}
