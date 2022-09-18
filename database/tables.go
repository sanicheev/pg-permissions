package database

func (c *Controller) ListTables() ([]Table, error) {
	tables := []Table{}

	err := c.Database.Select(
		&tables,
		`SELECT
			table_schema AS schema,
			table_name AS name
                FROM information_schema.tables`,
	)

	return tables, err
}

// return permission map
func (c *Controller) FetchTablePrivilege(user *string, table *string) (*PrivilegeDescriptor, error) {
	privilegeDescriptor := PrivilegeDescriptor{}

	err := c.Database.Get(
		&privilegeDescriptor,
		`SELECT
			has_table_privilege($1, $2, 'SELECT') AS select,
			has_table_privilege($1, $2, 'INSERT') AS insert,
			has_table_privilege($1, $2, 'UPDATE') AS update,
			has_table_privilege($1, $2, 'DELETE') AS delete,
			has_table_privilege($1, $2, 'REFERENCES') AS references,
			has_table_privilege($1, $2, 'TRIGGER') AS trigger,
			has_table_privilege($1, $2, 'TRUNCATE') AS truncate;`,
		*user,
		*table,
	)

	if err != nil {
		return nil, err
	}

	return &privilegeDescriptor, nil
}

func (c *Controller) TablePermissions(permissionType *string) []string {
	R := []string{
		"SELECT",
	}
	W := []string{
		"INSERT",
		"UPDATE",
		"DELETE",
		"REFERENCES",
		"TRIGGER",
		"TRUNCATE",
	}

	ALL := append(W, R...)

	if *permissionType == "W" {
		return W
	}

	return ALL
}
