package database

func (c *Controller) ListColumns() ([]Column, error) {
	columns := []Column{}

	err := c.Database.Select(
		&columns,
		`SELECT
			table_schema AS schema,
			table_name AS table,
			column_name AS name
                FROM information_schema.columns`,
	)

	return columns, err
}

func (c *Controller) FetchColumnPrivilege(user *string, table *string, column *string) (*PrivilegeDescriptor, error) {
	privilegeDescriptor := PrivilegeDescriptor{}

	err := c.Database.Get(
		&privilegeDescriptor,
		`SELECT
			has_column_privilege($1, $2, $3, 'SELECT') as select,
			has_column_privilege($1, $2, $3, 'INSERT') as insert,
			has_column_privilege($1, $2, $3, 'UPDATE') as update,
			has_column_privilege($1, $2, $3, 'REFERENCES') as references;`,
		*user,
		*table,
		*column,
	)

	if err != nil {
		return nil, err
	}

	return &privilegeDescriptor, nil
}

func (c *Controller) ColumnPermissions(permissionType *string) []string {
	R := []string{
		"SELECT",
	}

	W := []string{
		"INSERT",
		"UPDATE",
		"REFERENCES",
	}

	ALL := append(W, R...)

	if *permissionType == "W" {
		return W
	}

	return ALL
}
