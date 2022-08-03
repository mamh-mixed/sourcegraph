package cliutil

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/sourcegraph/log"
	"github.com/urfave/cli/v2"

	"github.com/sourcegraph/sourcegraph/internal/database"
	"github.com/sourcegraph/sourcegraph/internal/database/migration/schemas"
	"github.com/sourcegraph/sourcegraph/internal/oobmigration"
	"github.com/sourcegraph/sourcegraph/lib/output"
)

func RunOutOfBandMigrations(logger log.Logger, commandName string, runnerFactory RunnerFactory, outFactory OutputFactory, register func(db database.DB, outOfBandMigrationRunner *oobmigration.Runner) error) *cli.Command {
	action := makeAction(outFactory, func(ctx context.Context, cmd *cli.Context, out *output.Output) error {
		r, err := runnerFactory(ctx, schemas.SchemaNames)
		if err != nil {
			return err
		}
		db, err := extractDatabase(ctx, r)
		if err != nil {
			return err
		}

		store := oobmigration.NewStoreWithDB(db)
		outOfBandMigrationRunner := outOfBandMigrationRunnerWithStore(store)

		if err := register(db, outOfBandMigrationRunner); err != nil {
			return err
		}

		go outOfBandMigrationRunner.Start()
		defer outOfBandMigrationRunner.Stop()

		for range time.NewTicker(time.Second).C {
			migrations, err := store.List(ctx)
			if err != nil {
				return err
			}
			sort.Slice(migrations, func(i, j int) bool { return migrations[i].ID < migrations[j].ID })

			for _, m := range migrations {
				if !m.Complete() {
					fmt.Printf("> %d -> %.2f\n", m.ID, m.Progress*100)
				}
			}

			fmt.Printf("\n\n")
		}

		return nil
	})

	return &cli.Command{
		Name:        "run-out-of-band-migrations",
		Usage:       "TODO (Experimental)", // TODO
		Description: "",
		Action:      action,
		Flags:       []cli.Flag{},
	}
}
