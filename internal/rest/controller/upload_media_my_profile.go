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
	"github.com/netbill/restkit/render"
)

const operationGetMyProfileAvatarUploadMediaLink = "get_my_profile_avatar_upload_media_link"

func (c *Controller) CreateMyProfileUploadMediaLink(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetMyProfileAvatarUploadMediaLink)

	profile, media, err := c.modules.Profile.CreateUploadMediaLinks(r.Context(), scope.AccountActor(r))
	switch {
	case errors.Is(err, errx.ErrorProfileNotExists):
		log.Info("profile for user does not exist")
		render.ResponseError(w, problems.Unauthorized("profile for user does not exist"))
	case err != nil:
		log.WithError(err).Error("failed to open update profile session")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.UploadProfileMediaLinks(profile, media))
	}
}

const operationDeleteMyProfileUploadAvatar = "delete_my_profile_upload_avatar"

func (c *Controller) DeleteMyProfileUploadAvatar(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeleteMyProfileUploadAvatar)

	req, err := requests.DeleteUploadProfileAvatar(r)
	if err != nil {
		log.WithError(err).Info("invalid delete upload profile avatar request")
		render.ResponseError(w, problems.BadRequest(err)...)

		return
	}

	err = c.modules.Profile.DeleteUploadAvatar(
		r.Context(),
		scope.AccountActor(r),
		req.Data.Attributes.AvatarKey,
	)
	switch {
	case errors.Is(err, errx.ErrorProfileNotExists):
		log.Info("profile for user does not exist")
		render.ResponseError(w, problems.Unauthorized("profile for user does not exist"))
	case errors.Is(err, errx.ErrorProfileAvatarKeyIsInvalid):
		log.WithError(err).Info("avatar key is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"avatar": errors.New("avatar key is invalid"),
		})...)
	case err != nil:
		log.WithError(err).Error("failed to cancel update profile session")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, nil)
	}
}
