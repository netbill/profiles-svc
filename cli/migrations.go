package cli

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
)

//go:embed migrations/*.sql
var Migrations embed.FS

var migrations = &migrate.EmbedFileSystemMigrationSource{
	FileSystem: Migrations,
	Root:       "schema",
}

func openDB(ctx context.Context, url string) (*pgxpool.Pool, *sql.DB, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create pgx pool")
	}

	db := stdlib.OpenDBFromPool(pool)
	if err = db.PingContext(ctx); err != nil {
		db.Close()
		pool.Close()
		return nil, nil, errors.Wrap(err, "failed to ping database")
	}

	return pool, db, nil
}

func Migrate(ctx context.Context, url string, up bool, down bool) error {
	switch {
	case up && down:
		return fmt.Errorf("invalid migrate args: choose only one of --up or --down")
	case !up && !down:
		return fmt.Errorf("invalid migrate args: specify one of --up or --down")
	}

	pool, db, err := openDB(ctx, url)
	if err != nil {
		return err
	}
	defer db.Close()
	defer pool.Close()

	direction := migrate.Up
	if down {
		direction = migrate.Down
	}

	applied, err := migrate.ExecContext(ctx, db, "postgres", migrations, direction)
	if err != nil {
		if down {
			return errors.Wrap(err, "failed to apply migrations (down)")
		}
		return errors.Wrap(err, "failed to apply migrations (up)")
	}

	logrus.WithFields(logrus.Fields{
		"applied": applied,
		"dir":     map[bool]string{true: "down", false: "up"}[down],
	}).Info("migrations applied")

	return nil
}
