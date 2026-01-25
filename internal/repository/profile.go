package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/core/models"
	"github.com/netbill/profiles-svc/internal/core/modules/profile"
	"github.com/netbill/profiles-svc/internal/repository/pgdb"
	"github.com/netbill/restkit/pagi"
)

func (r Repository) InsertProfile(ctx context.Context, accountID uuid.UUID, username string) (models.Profile, error) {
	res, err := r.profilesQ(ctx).Insert(ctx, pgdb.InsertProfileParams{
		AccountID: accountID,
		Username:  username,
	})
	if err != nil {
		return models.Profile{}, fmt.Errorf(
			"failed to insert profile for account id %s, cause: %w", accountID, err,
		)
	}

	return res.ToModel(), nil
}

func (r Repository) GetProfileByAccountID(ctx context.Context, accountId uuid.UUID) (models.Profile, error) {
	row, err := r.profilesQ(ctx).FilterAccountID(accountId).Get(ctx)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Profile{}, errx.ErrorProfileNotFound.Raise(
			fmt.Errorf("failed to get profile by account id %s, cause: %w", accountId, err),
		)
	case err != nil:
		return models.Profile{}, fmt.Errorf(
			"failed to get profile by account id %s, cause: %w", accountId, err,
		)
	}

	return row.ToModel(), nil
}

func (r Repository) GetProfileByUsername(ctx context.Context, username string) (models.Profile, error) {
	row, err := r.profilesQ(ctx).FilterUsername(username).Get(ctx)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Profile{}, errx.ErrorProfileNotFound.Raise(
			fmt.Errorf("failed to get profile by username %s, cause: %w", username, err),
		)
	case err != nil:
		return models.Profile{}, fmt.Errorf(
			"failed to get profile by username %s, cause: %w", username, err,
		)
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
		if *input.Pseudonym == "" {
			q = q.UpdatePseudonym(pgtype.Text{
				String: "",
				Valid:  false,
			})
		} else {
			q = q.UpdatePseudonym(pgtype.Text{
				String: *input.Pseudonym,
				Valid:  true,
			})
		}
	}
	if input.Description != nil {
		if *input.Description == "" {
			q = q.UpdateDescription(pgtype.Text{
				String: "",
				Valid:  false,
			})
		} else {
			q = q.UpdateDescription(pgtype.Text{
				String: *input.Description,
				Valid:  true,
			})
		}
	}

	res, err := q.UpdateOne(ctx)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Profile{}, errx.ErrorProfileNotFound.Raise(
			fmt.Errorf("failed to update profile by account id %s, cause: %w", accountID, err),
		)
	case err != nil:
		return models.Profile{}, fmt.Errorf(
			"failed to update profile by account id %s, cause: %w", accountID, err,
		)
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
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Profile{}, errx.ErrorProfileNotFound.Raise(
			fmt.Errorf("failed to update profile username by account id %s, cause: %w", accountID, err),
		)
	case err != nil:
		return models.Profile{}, fmt.Errorf(
			"failed to update profile username by account id %s, cause: %w", accountID, err,
		)
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
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Profile{}, errx.ErrorProfileNotFound.Raise(
			fmt.Errorf("failed to update profile official by account id %s, cause: %w", accountID, err),
		)
	case err != nil:
		return models.Profile{}, fmt.Errorf(
			"failed to update profile official by account id %s, cause: %w", accountID, err,
		)
	}

	return res.ToModel(), nil
}

func (r Repository) UpdateProfileAvatar(
	ctx context.Context,
	accountID uuid.UUID,
	avatarURL string,
) (models.Profile, error) {
	res, err := r.profilesQ(ctx).
		FilterAccountID(accountID).
		UpdateAvatarURL(pgtype.Text{
			String: avatarURL,
			Valid:  true,
		}).
		UpdateOne(ctx)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Profile{}, errx.ErrorProfileNotFound.Raise(
			fmt.Errorf("failed to update profile avatar by account id %s, cause: %w", accountID, err),
		)
	case err != nil:
		return models.Profile{}, fmt.Errorf(
			"failed to update profile avatar by account id %s, cause: %w", accountID, err,
		)
	}

	return res.ToModel(), nil
}

func (r Repository) DeleteProfileAvatar(
	ctx context.Context,
	accountID uuid.UUID,
) (models.Profile, error) {
	res, err := r.profilesQ(ctx).
		FilterAccountID(accountID).
		UpdateAvatarURL(pgtype.Text{
			String: "",
			Valid:  false,
		}).
		UpdateOne(ctx)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Profile{}, errx.ErrorProfileNotFound.Raise(
			fmt.Errorf("failed to delete profile official by account id %s, cause: %w", accountID, err),
		)
	case err != nil:
		return models.Profile{}, fmt.Errorf(
			"failed to delete profile official by account id %s, cause: %w", accountID, err,
		)
	}

	return res.ToModel(), nil
}

func (r Repository) FilterProfilesByUsername(
	ctx context.Context,
	prefix string,
	offset uint,
	limit uint,
) (pagi.Page[[]models.Profile], error) {
	q := r.profilesQ(ctx).FilterLikeUsername(prefix)

	if limit == 0 {
		limit = 10
	}

	rows, err := q.Page(limit, offset).Select(ctx)
	if err != nil {
		return pagi.Page[[]models.Profile]{}, fmt.Errorf(
			"failed to filter profiles by username with prefix %s: %w", prefix, err,
		)
	}

	collection := make([]models.Profile, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, row.ToModel())
	}

	total, err := q.Count(ctx)
	if err != nil {
		return pagi.Page[[]models.Profile]{}, fmt.Errorf(
			"failed to count profiles by username with prefix %s: %w", prefix, err,
		)
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
		return pagi.Page[[]models.Profile]{}, fmt.Errorf(
			"failed to filter profiles: %w", err,
		)
	}

	collection := make([]models.Profile, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, row.ToModel())
	}

	total, err := q.Count(ctx)
	if err != nil {
		return pagi.Page[[]models.Profile]{}, fmt.Errorf(
			"failed to count profiles: %w", err,
		)
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
