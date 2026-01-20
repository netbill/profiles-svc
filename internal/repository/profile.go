package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/models"
	"github.com/netbill/profiles-svc/internal/core/modules/profile"
	"github.com/netbill/profiles-svc/internal/repository/pgdb"
	"github.com/netbill/restkit/pagi"
)

func (r Repository) CreateProfile(ctx context.Context, userID uuid.UUID, params profile.CreateParams) (models.Profile, error) {
	res, err := r.profilesQ(ctx).Insert(ctx, pgdb.InsertProfileParams{
		AccountID:   userID,
		Username:    params.Username,
		Pseudonym:   params.Pseudonym,
		Description: params.Description,
		Avatar:      params.Avatar,
		Official:    false,
	})
	if err != nil {
		return models.Profile{}, err
	}

	return res.ToModel(), nil
}

func (r Repository) GetProfileByAccountID(ctx context.Context, accountId uuid.UUID) (models.Profile, error) {
	row, err := r.profilesQ(ctx).FilterAccountID(accountId).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.Profile{}, nil
	case err != nil:
		return models.Profile{}, err
	}

	return row.ToModel(), nil
}

func (r Repository) GetProfileByUsername(ctx context.Context, username string) (models.Profile, error) {
	row, err := r.profilesQ(ctx).FilterUsername(username).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.Profile{}, nil
	case err != nil:
		return models.Profile{}, err
	}

	return row.ToModel(), nil
}

func (r Repository) UpdateProfile(
	ctx context.Context,
	accountID uuid.UUID,
	input profile.UpdateParams,
) (models.Profile, error) {
	q := r.profilesQ(ctx).FilterAccountID(accountID)

	if input.Pseudonym != nil {
		q = q.UpdatePseudonym(*input.Pseudonym)
	}
	if input.Description != nil {
		q = q.UpdateDescription(*input.Description)
	}
	if input.Avatar != nil {
		q = q.UpdateAvatar(*input.Avatar)
	}

	res, err := q.UpdateOne(ctx)
	if err != nil {
		return models.Profile{}, err
	}

	return res.ToModel(), nil
}

func (r Repository) UpdateProfileUsername(
	ctx context.Context,
	accountID uuid.UUID,
	username string,
) (models.Profile, error) {
	res, err := r.profilesQ(ctx).
		FilterAccountID(accountID).
		UpdateUsername(username).
		UpdateOne(ctx)
	if err != nil {
		return models.Profile{}, err
	}

	return res.ToModel(), nil
}

func (r Repository) UpdateProfileOfficial(
	ctx context.Context,
	accountID uuid.UUID,
	official bool,
) (models.Profile, error) {
	res, err := r.profilesQ(ctx).
		FilterAccountID(accountID).
		UpdateOfficial(official).
		UpdateOne(ctx)
	if err != nil {
		return models.Profile{}, err
	}

	return res.ToModel(), nil
}

func (r Repository) FilterProfilesByUsername(
	ctx context.Context,
	prefix string,
	offset uint,
	limit uint,
) (pagi.Page[[]models.Profile], error) {
	rows, err := r.profilesQ(ctx).
		FilterLikeUsername(prefix).
		Page(limit, offset).
		Select(ctx)
	if err != nil {
		return pagi.Page[[]models.Profile]{}, err
	}

	collection := make([]models.Profile, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, row.ToModel())
	}

	total, err := r.profilesQ(ctx).
		FilterLikeUsername(prefix).
		Count(ctx)
	if err != nil {
		return pagi.Page[[]models.Profile]{}, err
	}

	return pagi.Page[[]models.Profile]{
		Data:  collection,
		Page:  uint(offset/limit) + 1,
		Size:  uint(len(collection)),
		Total: total,
	}, nil
}

func (r Repository) FilterProfiles(
	ctx context.Context,
	params profile.FilterParams,
	limit, offset uint,
) (pagi.Page[[]models.Profile], error) {
	q := r.profilesQ(ctx)

	if params.PseudonymPrefix != nil {
		q = q.FilterLikePseudonym(*params.PseudonymPrefix)
	}
	if params.UsernamePrefix != nil {
		q = q.FilterLikeUsername(*params.UsernamePrefix)
	}

	if limit == 0 {
		limit = 10
	}

	rows, err := q.Page(limit, offset).Select(ctx)
	if err != nil {
		return pagi.Page[[]models.Profile]{}, err
	}

	collection := make([]models.Profile, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, row.ToModel())
	}

	total, err := q.Count(ctx)
	if err != nil {
		return pagi.Page[[]models.Profile]{}, err
	}

	return pagi.Page[[]models.Profile]{
		Data:  collection,
		Page:  uint(offset/limit) + 1,
		Size:  uint(len(collection)),
		Total: total,
	}, nil
}

func (r Repository) DeleteProfile(ctx context.Context, accountID uuid.UUID) error {
	return r.profilesQ(ctx).FilterAccountID(accountID).Delete(ctx)
}
