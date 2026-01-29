package repository

import (
	"context"
)

type Repository struct {
	profileSql ProfilesQ
	Transactioner
}

func New(Transaction Transactioner, profileSql ProfilesQ) Repository {
	return Repository{
		profileSql:    profileSql,
		Transactioner: Transaction,
	}
}

func (r Repository) profilesSqlQ() ProfilesQ {
	return r.profileSql.New()
}

type Transactioner interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
