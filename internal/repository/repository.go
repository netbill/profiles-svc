package repository

import (
	"context"
)

type Repository struct {
	ProfilesSQl ProfilesQ
	AccountsSql AccountsQ
	TombstonesSql
	TransactionSql
}

type TransactionSql interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
