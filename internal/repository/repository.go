package repository

import (
	"context"
)

type Repository struct {
	profilesQ ProfilesQ
	Transactioner
}

func New(Transaction Transactioner, profileSql ProfilesQ) *Repository {
	return &Repository{
		profilesQ:     profileSql,
		Transactioner: Transaction,
	}
}

type Transactioner interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
