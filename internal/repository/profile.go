package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/netbill/pagi"
	"github.com/netbill/profiles-svc/internal/domain/models"
	"github.com/netbill/profiles-svc/internal/domain/modules/profile"
	"github.com/netbill/profiles-svc/internal/repository/pgdb"
)

func (r Repository) CreateProfile(ctx context.Context, userID uuid.UUID, username string) (models.Profile, error) {
	res, err := r.profilesQ().Insert(ctx, pgdb.Profile{
		AccountID: userID,
		Username:  username,
	})
	if err != nil {
		return models.Profile{}, err
	}

	return res.ToModel(), nil
}

func (r Repository) GetProfileByAccountID(ctx context.Context, accountId uuid.UUID) (models.Profile, error) {
	row, err := r.profilesQ().FilterAccountID(accountId).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.Profile{}, nil
	case err != nil:
		return models.Profile{}, err
	}

	return row.ToModel(), nil
}

func (r Repository) GetProfileByUsername(ctx context.Context, username string) (models.Profile, error) {
	row, err := r.profilesQ().FilterUsername(username).Get(ctx)
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
	q := r.profilesQ().FilterAccountID(accountID)

	if input.Pseudonym != nil {
		q = q.UpdatePseudonym(input.Pseudonym)
	}
	if input.Description != nil {
		q = q.UpdateDescription(input.Description)
	}
	if input.Avatar != nil {
		q = q.UpdateAvatar(input.Avatar)
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
	res, err := r.profilesQ().
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
	res, err := r.profilesQ().
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
	rows, err := r.profilesQ().
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

	total, err := r.profilesQ().
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
	q := r.profilesQ()

	if params.PseudonymPrefix != nil {
		q = q.FilterLikePseudonym(*params.PseudonymPrefix)
	}
	if params.UsernamePrefix != nil {
		q = q.FilterLikeUsername(*params.UsernamePrefix)
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
	return r.profilesQ().FilterAccountID(accountID).Delete(ctx)
}
