package database

func (c *Controller) ListFunctions() ([]Function, error) {
	functions := []Function{}

	err := c.Database.Select(
		&functions,
		`SELECT
			r.routine_name AS name,
			r.routine_schema AS schema,
			r.routine_catalog AS catalog,
			pg_get_function_identity_arguments(pgproc.oid) AS arguments
		FROM information_schema.routines AS r
		INNER JOIN pg_catalog.pg_proc pgproc ON r.routine_name = pgproc.proname
		INNER JOIN pg_catalog.pg_namespace pgnamespace ON pgnamespace.oid = pgproc.pronamespace
		WHERE r.routine_type = 'FUNCTION' AND r.data_type != 'trigger';`,
	)

	return functions, err
}

func (c *Controller) FetchFunctionPrivilege(user *string, function *string) (*PrivilegeDescriptor, error) {
	privilegeDescriptor := PrivilegeDescriptor{}

	err := c.Database.Get(
		&privilegeDescriptor,
		`SELECT
			has_function_privilege($1, $2, 'EXECUTE') AS execute;`,
		*user,
		*function,
	)

	if err != nil {
		return nil, err
	}

	return &privilegeDescriptor, nil
}

func (c *Controller) FunctionPermissions(permissionType *string) []string {
	R := []string{
		"EXECUTE",
	}

	W := []string{}

	ALL := append(W, R...)

	if *permissionType == "W" {
		return W
	}

	return ALL
}
