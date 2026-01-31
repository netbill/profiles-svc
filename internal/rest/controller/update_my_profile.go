package controller

import (
	"errors"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/core/modules/profile"
	"github.com/netbill/profiles-svc/internal/rest/contexter"
	"github.com/netbill/restkit/problems"

	"github.com/netbill/profiles-svc/internal/rest/requests"
	"github.com/netbill/profiles-svc/internal/rest/responses"
)

func (c Controller) ConfirmUpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	initiator, err := contexter.AccountData(r.Context())
	if err != nil {
		c.log.WithError(err).Error("failed to get user from context")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.UpdateProfile(r)
	if err != nil {
		c.log.WithError(err).Errorf("invalid create profile request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	if req.Data.Id != initiator.GetAccountID() {
		c.log.WithError(err).Errorf("id in body and initiator id mismatch fir update My profile request")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"id": fmt.Errorf(
				"id in body: %s and initiator id: %s mismatch fir update My profile request",
				req.Data.Id,
				initiator.GetAccountID(),
			),
		})...)

		return
	}

	uploadData, err := contexter.UploadContentData(r.Context())
	if err != nil {
		c.log.WithError(err).Error("failed to get upload session id")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get upload session id"))

		return
	}

	res, err := c.domain.UpdateProfile(
		r.Context(),
		initiator.GetAccountID(),
		profile.UpdateParams{
			Pseudonym:   req.Data.Attributes.Pseudonym,
			Description: req.Data.Attributes.Description,
			Media: profile.UpdateMediaParams{
				UploadSessionID: uploadData.GetUploadSessionID(),
				DeleteAvatar:    req.Data.Attributes.DeleteAvatar,
			},
		},
	)
	if err != nil {
		c.log.WithError(err).Errorf("failed to update profile")
		switch {
		case errors.Is(err, errx.ErrorProfileNotFound):
			c.responser.RenderErr(w, problems.Unauthorized("profile for user does not exist"))
		case errors.Is(err, errx.ErrorProfileAvatarContentFormatIsNotAllowed),
			errors.Is(err, errx.ErrorProfileAvatarTooLarge),
			errors.Is(err, errx.ErrorProfileAvatarContentTypeIsNotAllowed):
			c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
				"avatar": fmt.Errorf(err.Error()),
			})...)
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}

		return
	}

	c.responser.Render(w, http.StatusOK, responses.Profile(res))
}
