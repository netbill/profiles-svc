package repository

import (
	"context"
)

type Repository struct {
	ProfilesQ ProfilesQ
	AccountsQ AccountsQ
	Transactioner
}

type Transactioner interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
