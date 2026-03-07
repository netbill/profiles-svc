package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/profiles-svc/internal/core/profile"
	"github.com/netbill/profiles-svc/internal/errx"
	"github.com/netbill/profiles-svc/internal/rest/requests"
	"github.com/netbill/profiles-svc/internal/rest/responses"
	"github.com/netbill/profiles-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationGetMyProfileAvatarUploadMediaLink = "get_my_profile_avatar_upload_media_link"

func (c *ProfileController) CreateUploadMediaLink(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetMyProfileAvatarUploadMediaLink)

	profile, media, err := c.profile.CreateUploadMediaLinks(r.Context(), scope.AccountActor(r))
	switch {
	case errors.Is(err, errx.ErrorProfileNotExists):
		log.Info("profile for user does not exist")
		render.ResponseError(w, problems.Unauthorized())
	case err != nil:
		log.WithError(err).Error("unexpected error")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.UploadProfileMediaLinks(r, profile, media))
	}
}

const operationDeleteMyProfileUploadAvatar = "delete_my_profile_upload_avatar"

func (c *ProfileController) DeleteUploadMedia(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeleteMyProfileUploadAvatar)

	req, err := requests.DeleteUploadProfileAvatar(r)
	if err != nil {
		log.WithError(err).Info("invalid delete upload profile avatar request")
		render.ResponseError(w, problems.BadRequest(err)...)

		return
	}

	log = log.With("target_avatar_id", req.Data.Id)

	err = c.profile.DeleteUploadMedia(
		r.Context(),
		scope.AccountActor(r),
		profile.DeleteUploadMediaParams{
			Avatar: req.Data.Attributes.AvatarKey,
		},
	)
	switch {
	case errors.Is(err, errx.ErrorProfileNotExists):
		log.WithError(err).Warn("profile for user does not exist")
		render.ResponseError(w, problems.Unauthorized())

	case err != nil:
		log.WithError(err).Error("unexpected error")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, nil)
	}
}
