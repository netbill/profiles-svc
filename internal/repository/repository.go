package repository

import (
	"context"
	"database/sql"

	"github.com/netbill/pgx"
	"github.com/netbill/profiles-svc/internal/repository/pgdb"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r Repository) profilesQ() pgdb.ProfilesQ {
	return pgdb.NewProfilesQ(r.db)
}

func (r Repository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return pgx.Transaction(r.db, ctx, fn)
}
