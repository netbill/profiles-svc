package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/profile"
	"github.com/netbill/profiles-svc/internal/errx"
	"github.com/netbill/profiles-svc/internal/models"
	"github.com/netbill/profiles-svc/internal/rest/requests"
	"github.com/netbill/profiles-svc/internal/rest/responses"
	"github.com/netbill/profiles-svc/internal/rest/scope"
	"github.com/netbill/restkit/pagi"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

type profileModule interface {
	GetList(
		ctx context.Context,
		params profile.FilterParams,
		limit, offset uint,
	) (pagi.Page[[]models.Profile], error)

	GetByID(ctx context.Context, accountID uuid.UUID) (models.Profile, error)
	GetByUsername(ctx context.Context, username string) (models.Profile, error)

	Update(
		ctx context.Context,
		actor models.AccountActor,
		params profile.UpdateParams,
	) (profile models.Profile, err error)

	CreateUploadMediaLinks(
		ctx context.Context,
		actor models.AccountActor,
	) (models.Profile, models.UploadProfileMediaLinks, error)

	DeleteUploadMedia(
		ctx context.Context,
		actor models.AccountActor,
		params profile.DeleteUploadMediaParams,
	) error
}

type ProfileController struct {
	profile profileModule
}

func New(profile profileModule) *ProfileController {
	return &ProfileController{
		profile: profile,
	}
}

const operationGetMyProfile = "get_my_profile"

func (c *ProfileController) GetMy(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetMyProfile)

	log = log.With("target_account_id", scope.AccountActor(r))

	res, err := c.profile.GetByID(r.Context(), scope.AccountActor(r))
	switch {
	case errors.Is(err, errx.ErrorProfileNotExists):
		log.WithError(err).Warn("profile for user does not exist")
		render.ResponseError(w, problems.Unauthorized())
	case err != nil:
		log.WithError(err).Error("unexpected error")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.Profile(r, res))
	}
}

const operationGetProfileByID = "get_profile_by_id"

func (c *ProfileController) GetByID(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetProfileByID)

	accountID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		log.WithError(err).Warn("invalid account id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"path": fmt.Errorf("invalid account id: %s", chi.URLParam(r, "account_id")),
		})...)
		return
	}

	log = log.With("target_account_id", accountID)

	res, err := c.profile.GetByID(r.Context(), accountID)
	switch {
	case errors.Is(err, errx.ErrorProfileNotExists):
		log.WithError(err).Warn("profile for user does not exist")
		render.ResponseError(w, problems.NotFound("profile for user does not exist"))
	case err != nil:
		log.WithError(err).Error("unexpected error")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.Profile(r, res))
	}
}

const operationGetProfileByUsername = "get_profile_by_username"

func (c *ProfileController) GetByUsername(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetProfileByUsername)

	username := chi.URLParam(r, "username")

	log = log.With("username", username)

	res, err := c.profile.GetByUsername(r.Context(), username)
	switch {
	case errors.Is(err, errx.ErrorProfileNotExists):
		log.WithError(err).Warn("profile for user does not exist")
		render.ResponseError(w, problems.NotFound("profile for user does not exist"))
	case err != nil:
		log.WithError(err).Error("unexpected error")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.Profile(r, res))
	}
}

const operationFilterProfiles = "filter_profiles"

func (c *ProfileController) Filter(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationFilterProfiles)

	q := r.URL.Query()
	limit, offset := pagi.GetPagination(r)

	filters := profile.FilterParams{}

	if text := strings.TrimSpace(q.Get("text")); text != "" {
		filters.Text = &text
	}

	res, err := c.profile.GetList(r.Context(), filters, limit, offset)
	switch {
	case err != nil:
		log.WithError(err).Error("unexpected error")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.ProfileCollection(r, res))
	}
}

const operationUpdateMyProfile = "update_my_profile"

func (c *ProfileController) UpdateMy(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationUpdateMyProfile)

	req, err := requests.UpdateProfile(r)
	if err != nil {
		log.WithError(err).Warn("invalid request body")
		render.ResponseError(w, problems.BadRequest(err)...)

		return
	}

	log = log.With("target_account_id", scope.AccountActor(r))

	res, err := c.profile.Update(r.Context(), scope.AccountActor(r), profile.UpdateParams{
		AvatarKey:   req.Data.Attributes.AvatarKey,
		Pseudonym:   req.Data.Attributes.Pseudonym,
		Description: req.Data.Attributes.Description,
	})
	switch {
	case errors.Is(err, errx.ErrorProfileNotExists):
		log.WithError(err).Warn("profile for user does not exist")
		render.ResponseError(w, problems.NotFound("profile for user does not exist"))
	case errors.Is(err, errx.ErrorProfileUploadedAvatarInvalid):
		log.WithError(err).Warn("avatar key is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"avatar": fmt.Errorf("avatar key is invalid"),
		})...)
	case errors.Is(err, errx.ErrorProfileUploadedAvatarInvalid):
		log.WithError(err).Warn("invalid avatar content")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"avatar": err,
		})...)
	case err != nil:
		log.WithError(err).Error("failed to update profile")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Debug("profile updated")
		render.Response(w, http.StatusOK, responses.Profile(r, res))
	}
}
