package controller

import (
	"errors"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/core/modules/profile"
	"github.com/netbill/profiles-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"

	"github.com/netbill/profiles-svc/internal/rest/requests"
	"github.com/netbill/profiles-svc/internal/rest/responses"
)

func (c *Controller) ConfirmUpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	req, err := requests.UpdateProfile(r)
	if err != nil {
		scope.Log(r).WithError(err).Errorf("invalid create profile request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	res, err := c.modules.Profile.ConfirmUpdateSession(
		r.Context(),
		scope.AccountAuthClaims(r).GetAccountID(),
		scope.UploadContentClaims(r).GetSessionID(),
		profile.UpdateParams{
			Pseudonym:   req.Data.Attributes.Pseudonym,
			Description: req.Data.Attributes.Description,
			Media: profile.UpdateMediaParams{
				DeleteAvatar: req.Data.Attributes.DeleteAvatar,
			},
		},
	)
	switch {
	case errors.Is(err, errx.ErrorProfileNotExists):
		scope.Log(r).Info("profile for user does not exist")
		c.responser.RenderErr(w, problems.Unauthorized("profile for user does not exist"))
	case errors.Is(err, errx.ErrorProfileAvatarContentIsInvalid):
		scope.Log(r).Info("avatar content is not valid for update profile")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"avatar": fmt.Errorf("avatar content is not valid for update profile"),
		})...)
	case err != nil:
		scope.Log(r).Errorf("failed to update profile")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusOK, responses.Profile(res))
	}
}
