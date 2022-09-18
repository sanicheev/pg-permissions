package cmd

import (
	"pg_permissions/database"
	"pg_permissions/postgres"
	"pg_permissions/types"

	"github.com/spf13/viper"
)

type PostgresClientCmd struct{}

func (pc *PostgresClientCmd) Run(v *viper.Viper) error {
	client := types.NewClient(v)

	dbController := database.NewController(client.Database)

	postgresController := postgres.NewController(
		dbController,
		client.Config.DatabasePermissions,
		client.Config.MaxRequests,
		client.Config.ReportHTML,
		client.Config.ReportJSON,
	)

	if err := postgresController.Run(); err != nil {
		return err
	}

	return nil
}
