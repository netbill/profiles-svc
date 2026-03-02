package controller

import (
	"errors"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/core/modules/profile"
	"github.com/netbill/profiles-svc/internal/rest/requests"
	"github.com/netbill/profiles-svc/internal/rest/responses"
	"github.com/netbill/profiles-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationUpdateMyProfile = "update_my_profile"

func (c *Controller) UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationUpdateMyProfile)

	req, err := requests.UpdateProfile(r)
	if err != nil {
		log.WithError(err).Warn("invalid request body")
		render.ResponseError(w, problems.BadRequest(err)...)

		return
	}

	log = log.With("target_account_id", scope.AccountActor(r))

	res, err := c.modules.Profile.Update(r.Context(), scope.AccountActor(r), profile.UpdateParams{
		AvatarKey:   req.Data.Attributes.AvatarKey,
		Pseudonym:   req.Data.Attributes.Pseudonym,
		Description: req.Data.Attributes.Description,
	})
	switch {
	case errors.Is(err, errx.ErrorProfileNotExists):
		log.WithError(err).Warn("profile for user does not exist")
		render.ResponseError(w, problems.NotFound("profile for user does not exist"))
	case errors.Is(err, errx.ErrorProfileAvatarKeyIsInvalid):
		log.WithError(err).Warn("avatar key is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"avatar": fmt.Errorf("avatar key is invalid"),
		})...)
	case errors.Is(err, errx.ErrorNoContentUploaded):
		log.WithError(err).Warn("no content uploaded for avatar")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"avatar": fmt.Errorf("no content uploaded for avatar"),
		})...)
	case errors.Is(err, errx.ErrorProfileAvatarContentIsExceedsMax):
		log.WithError(err).Warn("avatar content is exceeds max")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"avatar": err,
		})...)
	case errors.Is(err, errx.ErrorProfileAvatarResolutionIsInvalid):
		log.WithError(err).Warn("avatar resolution is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"avatar": err,
		})...)
	case errors.Is(err, errx.ErrorProfileAvatarFormatIsNotAllowed):
		log.WithError(err).Warn("avatar format is not allowed")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"avatar": err,
		})...)
	case err != nil:
		log.WithError(err).Error("failed to update profile")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Debug("profile updated")
		render.Response(w, http.StatusOK, responses.Profile(res))
	}
}
