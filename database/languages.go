package database

func (c *Controller) ListLanguages() ([]Language, error) {
	languages := []Language{}

	err := c.Database.Select(
		&languages,
		`SELECT
			lanname AS language
		FROM pg_catalog.pg_language;`,
	)

	return languages, err
}

func (c *Controller) FetchLanguagePrivilege(user *string, language *string) (*PrivilegeDescriptor, error) {
	privilegeDescriptor := PrivilegeDescriptor{}

	err := c.Database.Get(
		&privilegeDescriptor,
		`SELECT
			has_language_privilege($1, $2, 'USAGE') AS usage;`,
		*user,
		*language,
	)

	if err != nil {
		return nil, err
	}

	return &privilegeDescriptor, nil
}

func (c *Controller) LanguagePermissions(permissionType *string) []string {
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
