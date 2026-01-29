package pg

import (
	"context"

	"github.com/netbill/pgdbx"
	"github.com/netbill/profiles-svc/internal/repository"
)

type transaction struct {
	db *pgdbx.DB
}

func NewTransaction(db *pgdbx.DB) repository.Transactioner {
	return &transaction{
		db: db,
	}
}

func (q *transaction) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return q.db.Transaction(ctx, fn)
}
