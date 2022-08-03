package codeintel

import (
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sourcegraph/log"

	"github.com/sourcegraph/sourcegraph/internal/codeintel/stores"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/stores/dbstore"
	dbmigrations "github.com/sourcegraph/sourcegraph/internal/codeintel/stores/dbstore/migration"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/stores/gitserver"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/stores/lsifstore"
	lsifmigrations "github.com/sourcegraph/sourcegraph/internal/codeintel/stores/lsifstore/migration"
	"github.com/sourcegraph/sourcegraph/internal/database"
	"github.com/sourcegraph/sourcegraph/internal/observation"
	"github.com/sourcegraph/sourcegraph/internal/oobmigration"
	"github.com/sourcegraph/sourcegraph/internal/trace"
)

// RegisterMigrations registers all code intel related out-of-band migration instances that should run for the current version of Sourcegraph.
func RegisterMigrations(db database.DB, outOfBandMigrationRunner *oobmigration.Runner) error {
	if err := config.Validate(); err != nil {
		return err
	}

	observationContext := &observation.Context{
		Logger:     log.Scoped("store", "codeintel db store"), // TODO
		Tracer:     &trace.Tracer{Tracer: opentracing.GlobalTracer()},
		Registerer: prometheus.DefaultRegisterer,
	}
	dbStore := dbstore.NewWithDB(db, observationContext)
	lsifStore := lsifstore.NewStore(stores.NewCodeIntelDBWith(db), nil, observationContext)
	gitserverClient := gitserver.New(db, dbStore, observationContext)

	if err := outOfBandMigrationRunner.Register(
		lsifmigrations.DiagnosticsCountMigrationID, // 1
		lsifmigrations.NewDiagnosticsCountMigrator(lsifStore, config.DiagnosticsCountMigrationBatchSize),
		oobmigration.MigratorOptions{Interval: config.DiagnosticsCountMigrationBatchInterval},
	); err != nil {
		return err
	}

	if err := outOfBandMigrationRunner.Register(
		lsifmigrations.DefinitionsCountMigrationID, // 4
		lsifmigrations.NewLocationsCountMigrator(lsifStore, "lsif_data_definitions", config.DefinitionsCountMigrationBatchSize),
		oobmigration.MigratorOptions{Interval: config.DefinitionsCountMigrationBatchInterval},
	); err != nil {
		return err
	}

	if err := outOfBandMigrationRunner.Register(
		lsifmigrations.ReferencesCountMigrationID, // 5
		lsifmigrations.NewLocationsCountMigrator(lsifStore, "lsif_data_references", config.ReferencesCountMigrationBatchSize),
		oobmigration.MigratorOptions{Interval: config.ReferencesCountMigrationBatchInterval},
	); err != nil {
		return err
	}

	if err := outOfBandMigrationRunner.Register(
		lsifmigrations.DocumentColumnSplitMigrationID, // 7
		lsifmigrations.NewDocumentColumnSplitMigrator(lsifStore, config.DocumentColumnSplitMigrationBatchSize),
		oobmigration.MigratorOptions{Interval: config.DocumentColumnSplitMigrationBatchInterval},
	); err != nil {
		return err
	}

	if err := outOfBandMigrationRunner.Register(
		lsifmigrations.APIDocsSearchMigrationID, // 12
		lsifmigrations.NewAPIDocsSearchMigrator(config.APIDocsSearchMigrationBatchSize),
		oobmigration.MigratorOptions{Interval: config.APIDocsSearchMigrationBatchInterval},
	); err != nil {
		return err
	}

	if err := outOfBandMigrationRunner.Register(
		dbmigrations.CommittedAtMigrationID, // 8
		dbmigrations.NewCommittedAtMigrator(dbStore, gitserverClient, config.CommittedAtMigrationBatchSize),
		oobmigration.MigratorOptions{Interval: config.CommittedAtMigrationBatchInterval},
	); err != nil {
		return err
	}

	if err := outOfBandMigrationRunner.Register(
		dbmigrations.ReferenceCountMigrationID, // 11
		dbmigrations.NewReferenceCountMigrator(dbStore, config.ReferenceCountMigrationBatchSize),
		oobmigration.MigratorOptions{Interval: config.ReferenceCountMigrationBatchInterval},
	); err != nil {
		return err
	}

	return nil
}
