package profile

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/errx"
	"github.com/netbill/profiles-svc/internal/models"

	"github.com/netbill/restkit/pagi"
)

type profileMessenger interface {
	WriteProfileUpdated(ctx context.Context, profile models.Profile) error
}

type usernameValidator interface {
	Validate(username string) error
}

type Service struct {
	repo      profileRepo
	tx        transaction
	messenger profileMessenger
	bucket    media
	username  usernameValidator
}

type ServiceDeps struct {
	Repo      profileRepo
	Tx        transaction
	Messenger profileMessenger
	Bucket    media
	Username  usernameValidator
}

func NewProfileModule(deps ServiceDeps) *Service {
	return &Service{
		repo:      deps.Repo,
		tx:        deps.Tx,
		messenger: deps.Messenger,
		bucket:    deps.Bucket,
		username:  deps.Username,
	}
}

func (s *Service) GetByID(
	ctx context.Context,
	accountID uuid.UUID,
) (models.Profile, error) {
	res, err := s.repo.GetByID(ctx, accountID)
	if err != nil {
		return models.Profile{}, err
	}

	return res, nil
}

func (s *Service) GetByUsername(
	ctx context.Context,
	username string,
) (models.Profile, error) {
	return s.repo.GetByUsername(ctx, username)
}

type FilterParams struct {
	Text *string
}

func (s *Service) GetList(
	ctx context.Context,
	params FilterParams,
	limit, offset uint,
) (pagi.Page[[]models.Profile], error) {
	collection, err := s.repo.Filter(ctx, params, limit, offset)
	if err != nil {
		return pagi.Page[[]models.Profile]{}, err
	}

	return collection, nil
}

type UpdateParams struct {
	AvatarKey   *string
	Pseudonym   *string
	Description *string
}

func (p UpdateParams) HasChanges(model models.Profile) bool {
	return !ptrEqual(p.AvatarKey, model.AvatarKey) ||
		!ptrEqual(p.Pseudonym, model.Pseudonym) ||
		!ptrEqual(p.Description, model.Description)
}

func ptrEqual[T comparable](a, b *T) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func (s *Service) Update(
	ctx context.Context,
	actor models.AccountActor,
	params UpdateParams,
) (profile models.Profile, err error) {
	profile, err = s.repo.GetByID(ctx, actor)
	if err != nil {
		return models.Profile{}, err
	}

	if !params.HasChanges(profile) {
		return profile, nil
	}

	switch {
	case params.AvatarKey != nil && *params.AvatarKey == "" && profile.AvatarKey != nil:
		if err := s.bucket.DeleteProfileAvatar(ctx, actor, *profile.AvatarKey); err != nil {
			return models.Profile{}, fmt.Errorf("failed to delete profile avatar: %w", err)
		}
		params.AvatarKey = nil
	case params.AvatarKey != nil:
		avatarKey, err := s.bucket.UpdateProfileAvatar(ctx, actor, *params.AvatarKey)
		if err != nil {
			return models.Profile{}, fmt.Errorf("failed to update profile avatar: %w", err)
		}
		params.AvatarKey = &avatarKey
	}

	if err = s.tx.Transaction(ctx, func(ctx context.Context) error {
		profile, err = s.repo.Update(ctx, actor, params)
		if err != nil {
			return err
		}

		if err = s.messenger.WriteProfileUpdated(ctx, profile); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return models.Profile{}, err
	}

	return profile, nil
}

func (s *Service) UpdateUsername(
	ctx context.Context,
	actor models.AccountActor,
	newUsername string,
) (models.Profile, error) {
	if err := s.username.Validate(newUsername); err != nil {
		return models.Profile{}, err
	}

	current, err := s.repo.GetByID(ctx, actor)
	if err != nil {
		return models.Profile{}, err
	}

	if current.Username == newUsername {
		return current, nil
	}

	unavalible, err := s.repo.ExistByUsername(ctx, newUsername)
	if err != nil {
		return models.Profile{}, err
	}
	if unavalible {
		return models.Profile{}, errx.ErrorProfileUsernameTaken
	}

	var profile models.Profile
	if err = s.tx.Transaction(ctx, func(ctx context.Context) error {
		profile, err = s.repo.UpdateUsername(ctx, actor, newUsername)
		if err != nil {
			return err
		}

		return s.messenger.WriteProfileUpdated(ctx, profile)
	}); err != nil {
		return models.Profile{}, err
	}

	return profile, nil
}
