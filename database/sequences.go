package database

func (c *Controller) ListSequences() ([]Sequence, error) {
	sequences := []Sequence{}

	err := c.Database.Select(
		&sequences,
		`SELECT
			sequence_schema AS schema,
			sequence_name AS name
                FROM information_schema.sequences`,
	)

	return sequences, err
}

func (c *Controller) FetchSequencePrivilege(user *string, sequence *string) (*PrivilegeDescriptor, error) {
	privilegeDescriptor := PrivilegeDescriptor{}

	err := c.Database.Get(
		&privilegeDescriptor,
		`SELECT
			has_sequence_privilege($1, $2, 'USAGE') AS usage,
			has_sequence_privilege($1, $2, 'SELECT') AS select,
			has_sequence_privilege($1, $2, 'UPDATE') AS update`,
		*user,
		*sequence,
	)

	if err != nil {
		return nil, err
	}

	return &privilegeDescriptor, nil
}

func (c *Controller) SequencePermissions(permissionType *string) []string {
	R := []string{
		"USAGE",
		"SELECT",
	}
	W := []string{
		"UPDATE",
	}

	ALL := append(W, R...)

	if *permissionType == "W" {
		return W
	}

	return ALL
}
