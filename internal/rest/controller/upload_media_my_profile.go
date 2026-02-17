package controller

import (
	"errors"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/rest/requests"
	"github.com/netbill/profiles-svc/internal/rest/responses"
	"github.com/netbill/profiles-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationGetMyProfileAvatarUploadMediaLink = "get_my_profile_avatar_upload_media_link"

func (c *Controller) CreateMyProfileUploadMediaLink(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetMyProfileAvatarUploadMediaLink)

	profile, media, err := c.modules.Profile.CreateProfileUploadMediaLinks(r.Context(), scope.AccountActor(r))
	switch {
	case errors.Is(err, errx.ErrorProfileNotExists):
		log.Info("profile for user does not exist")
		c.responser.RenderErr(w, problems.Unauthorized("profile for user does not exist"))
	case err != nil:
		log.WithError(err).Error("failed to open update profile session")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		log.Debug("update profile session opened")
		c.responser.Render(w, http.StatusOK, responses.UploadProfileMediaLinks(profile, media))
	}
}

const operationDeleteMyProfileUploadAvatar = "delete_my_profile_upload_avatar"

func (c *Controller) DeleteMyProfileUploadAvatar(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeleteMyProfileUploadAvatar)

	req, err := requests.DeleteUploadProfileAvatar(r)
	if err != nil {
		log.WithError(err).Info("invalid delete upload profile avatar request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	err = c.modules.Profile.DeleteProfileUploadAvatar(
		r.Context(),
		scope.AccountActor(r),
		req.Data.Attributes.AvatarKey,
	)
	switch {
	case errors.Is(err, errx.ErrorProfileNotExists):
		log.Info("profile for user does not exist")
		c.responser.RenderErr(w, problems.Unauthorized("profile for user does not exist"))
	case errors.Is(err, errx.ErrorProfileAvatarKeyIsInvalid):
		log.WithError(err).Info("avatar key is invalid")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"avatar": errors.New("avatar key is invalid"),
		})...)
	case err != nil:
		log.WithError(err).Error("failed to cancel update profile session")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		log.Debug("profile update session cancelled")
		c.responser.Status(w, http.StatusOK)
	}
}
