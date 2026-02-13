package cli

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	eventpg "github.com/netbill/eventbox/pg"
	"github.com/netbill/logium"
	"github.com/netbill/pgdbx"
)

func CleanupInboxFailed(ctx context.Context, cfg Config, log *logium.Entry) error {
	pool, err := pgxpool.New(ctx, cfg.Database.SQL.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pool.Close()

	db := pgdbx.NewDB(pool)

	inboxCleaner := eventpg.NewInboxCleaner(db)

	err = inboxCleaner.CleanInboxFailed(ctx)
	if err != nil {
		log.WithError(err).Error("failed to clean inbox failed")
		return err
	}

	log.Info("inbox failed cleaned successfully")
	return nil
}

func CleanupInboxProcessing(ctx context.Context, cfg Config, log *logium.Entry, processIDs ...string) error {
	pool, err := pgxpool.New(ctx, cfg.Database.SQL.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pool.Close()

	db := pgdbx.NewDB(pool)

	inboxCleaner := eventpg.NewInboxCleaner(db)

	err = inboxCleaner.CleanInboxProcessing(ctx, processIDs...)
	if err != nil {
		log.WithError(err).Error("failed to clean inbox processing")
		return err
	}

	log.Info("inbox processing cleaned successfully")
	return nil
}

func CleanupOutboxFailed(ctx context.Context, cfg Config, log *logium.Entry) error {
	pool, err := pgxpool.New(ctx, cfg.Database.SQL.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	db := pgdbx.NewDB(pool)

	OutboxCleaner := eventpg.NewOutboxCleaner(db)

	err = OutboxCleaner.CleanOutboxFailed(ctx)
	if err != nil {
		log.WithError(err).Error("failed to clean Outbox failed")
		return err
	}

	log.Info("Outbox failed cleaned successfully")
	return nil
}

func CleanupOutboxProcessing(ctx context.Context, cfg Config, log *logium.Entry, processIDs ...string) error {
	pool, err := pgxpool.New(ctx, cfg.Database.SQL.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	db := pgdbx.NewDB(pool)

	OutboxCleaner := eventpg.NewOutboxCleaner(db)

	err = OutboxCleaner.CleanOutboxProcessing(ctx, processIDs...)
	if err != nil {
		log.WithError(err).Error("failed to clean Outbox processing")
		return err
	}

	log.Info("Outbox processing cleaned successfully")
	return nil
}
