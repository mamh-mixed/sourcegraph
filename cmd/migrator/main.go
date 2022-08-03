package main

import (
	"github.com/sourcegraph/log"

	"github.com/sourcegraph/sourcegraph/cmd/migrator/shared"
	"github.com/sourcegraph/sourcegraph/internal/env"
	ossmigrations "github.com/sourcegraph/sourcegraph/internal/oobmigration/migrations"
	"github.com/sourcegraph/sourcegraph/internal/version"
)

func main() {
	liblog := log.Init(log.Resource{
		Name:    env.MyName,
		Version: version.Version(),
	})
	defer liblog.Sync()

	logger := log.Scoped("migrator", "migrator oss edition")

	if err := shared.Start(logger, ossmigrations.RegisterOSSMigrations); err != nil {
		logger.Fatal(err.Error())
	}
}
