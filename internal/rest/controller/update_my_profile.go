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
)

const operationUpdateMyProfile = "update_my_profile"

func (c *Controller) UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationUpdateMyProfile)

	req, err := requests.UpdateProfile(r)
	if err != nil {
		log.WithError(err).Info("invalid request body")
		c.responser.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	res, err := c.modules.Profile.Update(r.Context(), scope.AccountActor(r), profile.UpdateParams{
		AvatarKey:   req.Data.Attributes.AvatarKey,
		Pseudonym:   req.Data.Attributes.Pseudonym,
		Description: req.Data.Attributes.Description,
	})
	switch {
	case errors.Is(err, errx.ErrorProfileNotExists):
		log.WithError(err).Info("account not found by credentials")
		c.responser.RenderErr(w, problems.NotFound("profile for user does not exist"))
	case errors.Is(err, errx.ErrorProfileAvatarKeyIsInvalid):
		log.WithError(err).Info("avatar key is invalid")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"avatar": fmt.Errorf("avatar key is invalid"),
		})...)
	case errors.Is(err, errx.ErrorNoContentUploaded):
		log.WithError(err).Info("no content uploaded for avatar")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"avatar": fmt.Errorf("no content uploaded for avatar"),
		})...)
	case errors.Is(err, errx.ErrorProfileAvatarContentIsExceedsMax):
		log.WithError(err).Info("avatar content is exceeds max")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"avatar": err,
		})...)
	case errors.Is(err, errx.ErrorProfileAvatarResolutionIsInvalid):
		log.WithError(err).Info("avatar resolution is invalid")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"avatar": err,
		})...)
	case errors.Is(err, errx.ErrorProfileAvatarFormatIsNotAllowed):
		log.WithError(err).Info("avatar format is not allowed")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"avatar": err,
		})...)
	case err != nil:
		log.WithError(err).Error("failed to update profile")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		log.Debug("profile updated")
		c.responser.Render(w, http.StatusOK, responses.Profile(res))
	}
}
