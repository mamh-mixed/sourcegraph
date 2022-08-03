package main

import (
	"github.com/sourcegraph/log"

	"github.com/sourcegraph/sourcegraph/cmd/migrator/shared"
	enterprisemigrations "github.com/sourcegraph/sourcegraph/enterprise/internal/oobmigration/migrations"
	"github.com/sourcegraph/sourcegraph/internal/database"
	"github.com/sourcegraph/sourcegraph/internal/env"
	"github.com/sourcegraph/sourcegraph/internal/oobmigration"
	ossmigrations "github.com/sourcegraph/sourcegraph/internal/oobmigration/migrations"
	"github.com/sourcegraph/sourcegraph/internal/version"
)

func main() {
	liblog := log.Init(log.Resource{
		Name:    env.MyName,
		Version: version.Version(),
	})
	defer liblog.Sync()

	logger := log.Scoped("migrator", "migrator enterprise edition")

	if err := shared.Start(logger, register); err != nil {
		logger.Fatal(err.Error())
	}
}

func init() {
	oobmigration.ReturnEnterpriseMigrations = true
}

func register(db database.DB, outOfBandMigrationRunner *oobmigration.Runner) error {
	if err := ossmigrations.RegisterOSSMigrations(db, outOfBandMigrationRunner); err != nil {
		return err
	}
	if err := enterprisemigrations.RegisterEnterpriseMigrations(db, outOfBandMigrationRunner); err != nil {
		return err
	}

	return nil
}
