package profile

import (
	"context"

	"github.com/netbill/profiles-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
)

type FilterParams struct {
	UsernamePrefix  *string
	PseudonymPrefix *string
	Verified        *bool
}

func (s Service) FilterProfile(
	ctx context.Context,
	params FilterParams,
	limit, offset uint,
) (pagi.Page[[]models.Profile], error) {
	collection, err := s.repo.FilterProfiles(ctx, params, limit, offset)
	if err != nil {
		return pagi.Page[[]models.Profile]{}, err
	}

	return collection, nil
}
