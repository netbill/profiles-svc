package profile

import (
	"context"
	"fmt"

	"github.com/netbill/pagi"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/core/models"
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
		return pagi.Page[[]models.Profile]{}, errx.ErrorInternal.Raise(
			fmt.Errorf("getting profile with username '%s': %w", *params.UsernamePrefix, err),
		)
	}

	return collection, nil
}
