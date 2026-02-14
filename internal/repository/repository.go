package repository

import (
	"context"
)

type Repository struct {
	profilesQ ProfilesQ
	accountsQ AccountsQ
	Transactioner
}

func New(Transaction Transactioner, accountsSql AccountsQ, profileSql ProfilesQ) *Repository {
	return &Repository{
		profilesQ:     profileSql,
		accountsQ:     accountsSql,
		Transactioner: Transaction,
	}
}

type Transactioner interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
