package database

func (c *Controller) ListSchemas() ([]Schema, error) {
	schemas := []Schema{}

	err := c.Database.Select(
		&schemas,
		`SELECT
			schema_name AS name,
			schema_owner AS owner
                FROM information_schema.schemata`,
	)

	return schemas, err
}

func (c *Controller) FetchSchemaPrivilege(user *string, schema *string) (*PrivilegeDescriptor, error) {
	privilegeDescriptor := PrivilegeDescriptor{}

	err := c.Database.Get(
		&privilegeDescriptor,
		`SELECT
			has_schema_privilege($1, $2, 'USAGE') AS usage,
			has_schema_privilege($1, $2, 'CREATE') AS create;`,
		*user,
		*schema,
	)

	if err != nil {
		return nil, err
	}

	return &privilegeDescriptor, nil
}

func (c *Controller) SchemaPermissions(permissionType *string) []string {
	R := []string{
		"USAGE",
	}
	W := []string{
		"CREATE",
	}

	ALL := append(W, R...)

	if *permissionType == "W" {
		return W
	}

	return ALL
}
