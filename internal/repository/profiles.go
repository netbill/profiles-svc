package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/core/models"
	"github.com/netbill/profiles-svc/internal/core/modules/profile"
	"github.com/netbill/restkit/pagi"
)

type ProfileRow struct {
	AccountID   uuid.UUID `db:"account_id"`
	Username    string    `db:"username"`
	Official    bool      `db:"official"`
	Pseudonym   *string   `db:"pseudonym,omitempty"`
	Description *string   `db:"description,omitempty"`
	Avatar      *string   `db:"avatar,omitempty"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (p ProfileRow) IsNil() bool {
	return p.AccountID == uuid.Nil
}

func (p ProfileRow) ToModel() models.Profile {
	return models.Profile{
		AccountID:   p.AccountID,
		Username:    p.Username,
		Official:    p.Official,
		Pseudonym:   p.Pseudonym,
		Description: p.Description,
		Avatar:      p.Avatar,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

type ProfilesQ interface {
	New() ProfilesQ
	Insert(ctx context.Context, input ProfileRow) (ProfileRow, error)

	Get(ctx context.Context) (ProfileRow, error)
	Select(ctx context.Context) ([]ProfileRow, error)

	UpdateMany(ctx context.Context) (int64, error)
	UpdateOne(ctx context.Context) (ProfileRow, error)

	UpdateUsername(username string) ProfilesQ
	UpdateOfficial(official bool) ProfilesQ
	UpdatePseudonym(v *string) ProfilesQ
	UpdateDescription(v *string) ProfilesQ
	UpdateAvatar(v *string) ProfilesQ

	Delete(ctx context.Context) error

	FilterAccountID(accountID ...uuid.UUID) ProfilesQ
	FilterUsername(username string) ProfilesQ
	FilterOfficial(official bool) ProfilesQ
	FilterBestMatch(term string) ProfilesQ

	Count(ctx context.Context) (uint, error)
	Page(limit, offset uint) ProfilesQ
}

func (r *Repository) InsertProfile(ctx context.Context, accountID uuid.UUID, username string) (models.Profile, error) {
	res, err := r.profilesQ.Insert(ctx, ProfileRow{
		AccountID: accountID,
		Username:  username,
		Official:  false,
	})
	if err != nil {
		return models.Profile{}, fmt.Errorf(
			"failed to insert profile for account id %s, cause: %w", accountID, err,
		)
	}

	return res.ToModel(), nil
}

func (r *Repository) GetProfileByAccountID(ctx context.Context, accountID uuid.UUID) (models.Profile, error) {
	row, err := r.profilesQ.FilterAccountID(accountID).Get(ctx)
	switch {
	case err != nil:
		return models.Profile{}, fmt.Errorf(
			"failed to get profile by account id %s, cause: %w", accountID, err,
		)
	case row.IsNil():
		return models.Profile{}, errx.ErrorProfileNotExists.Raise(
			fmt.Errorf("profile by account id %s: profile not found", accountID),
		)
	}

	return row.ToModel(), nil
}

func (r *Repository) GetProfileByUsername(ctx context.Context, username string) (models.Profile, error) {
	row, err := r.profilesQ.FilterUsername(username).Get(ctx)
	switch {
	case err != nil:
		return models.Profile{}, fmt.Errorf(
			"failed to get profile by username %s, cause: %w", username, err,
		)
	case row.IsNil():
		return models.Profile{}, errx.ErrorProfileNotExists.Raise(
			fmt.Errorf("failed to get profile by username %s, cause: %w", username, err),
		)
	}

	return row.ToModel(), nil
}

func (r *Repository) UpdateProfile(
	ctx context.Context,
	accountID uuid.UUID,
	input profile.UpdateParams,
) (models.Profile, error) {
	q := r.profilesQ.
		FilterAccountID(accountID).
		UpdatePseudonym(input.Pseudonym).
		UpdateDescription(input.Description).
		UpdateAvatar(input.GetUpdatedAvatar())

	row, err := q.UpdateOne(ctx)
	switch {
	case err != nil:
		return models.Profile{}, fmt.Errorf(
			"failed to update profile by account id %s, cause: %w", accountID, err,
		)
	case row.IsNil():
		return models.Profile{}, errx.ErrorProfileNotExists.Raise(
			fmt.Errorf("failed to update profile by account id %s, cause: %w", accountID, err),
		)
	}

	return row.ToModel(), nil
}

func (r *Repository) UpdateProfileUsername(
	ctx context.Context,
	accountID uuid.UUID,
	username string,
) (models.Profile, error) {
	row, err := r.profilesQ.
		FilterAccountID(accountID).
		UpdateUsername(username).
		UpdateOne(ctx)
	switch {
	case err != nil:
		return models.Profile{}, fmt.Errorf(
			"failed to update profile username by account id %s, cause: %w", accountID, err,
		)
	case row.IsNil():
		return models.Profile{}, errx.ErrorProfileNotExists.Raise(
			fmt.Errorf("failed to update profile username by account id %s, cause: %w", accountID, err),
		)
	}

	return row.ToModel(), nil
}

func (r *Repository) UpdateProfileOfficial(
	ctx context.Context,
	accountID uuid.UUID,
	official bool,
) (models.Profile, error) {
	row, err := r.profilesQ.
		FilterAccountID(accountID).
		UpdateOfficial(official).
		UpdateOne(ctx)
	switch {
	case err != nil:
		return models.Profile{}, fmt.Errorf(
			"failed to update profile official by account id %s, cause: %w", accountID, err,
		)
	case row.IsNil():
		return models.Profile{}, errx.ErrorProfileNotExists.Raise(
			fmt.Errorf("failed to update profile official by account id %s, cause: %w", accountID, err),
		)
	}

	return row.ToModel(), nil
}

func (r *Repository) UpdateProfileAvatar(
	ctx context.Context,
	accountID uuid.UUID,
	avatarURL string,
) (models.Profile, error) {
	row, err := r.profilesQ.
		FilterAccountID(accountID).
		UpdateAvatar(&avatarURL).
		UpdateOne(ctx)
	switch {
	case err != nil:
		return models.Profile{}, fmt.Errorf(
			"failed to update profile avatar by account id %s, cause: %w", accountID, err,
		)
	case row.IsNil():
		return models.Profile{}, errx.ErrorProfileNotExists.Raise(
			fmt.Errorf("failed to update profile avatar by account id %s, cause: %w", accountID, err),
		)
	}

	return row.ToModel(), nil
}

func (r *Repository) DeleteProfileAvatar(
	ctx context.Context,
	accountID uuid.UUID,
) (models.Profile, error) {
	row, err := r.profilesQ.
		FilterAccountID(accountID).
		UpdateAvatar(nil).
		UpdateOne(ctx)
	switch {
	case err != nil:
		return models.Profile{}, fmt.Errorf(
			"failed to delete profile official by account id %s, cause: %w", accountID, err,
		)
	case row.IsNil():
		return models.Profile{}, errx.ErrorProfileNotExists.Raise(
			fmt.Errorf("failed to delete profile official by account id %s, cause: %w", accountID, err),
		)
	}

	return row.ToModel(), nil
}

func (r *Repository) FilterProfiles(
	ctx context.Context,
	params profile.FilterParams,
	limit, offset uint,
) (pagi.Page[[]models.Profile], error) {
	q := r.profilesQ

	if params.Text != nil {
		q = q.FilterBestMatch(*params.Text)
	}
	if params.Official != nil {
		q = q.FilterOfficial(*params.Official)
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

func (r *Repository) DeleteProfile(ctx context.Context, accountID uuid.UUID) error {
	return r.profilesQ.FilterAccountID(accountID).Delete(ctx)
}
