package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/netbill/pgxtx"
	"github.com/netbill/profiles-svc/internal/repository/pgdb"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) Repository {
	return Repository{pool: pool}
}

func (r Repository) profilesQ(ctx context.Context) pgdb.ProfilesQ {
	return pgdb.NewProfilesQ(pgxtx.Exec(r.pool, ctx))
}

func (r Repository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return pgxtx.Transaction(r.pool, ctx, fn)
}
