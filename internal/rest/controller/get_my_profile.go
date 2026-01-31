package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/rest/contexter"
	"github.com/netbill/profiles-svc/internal/rest/responses"
	"github.com/netbill/restkit/problems"
)

func (c Controller) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	initiator, err := contexter.AccountData(r.Context())
	if err != nil {
		c.log.WithError(err).Error("failed to get account from context")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get account from context"))

		return
	}

	res, err := c.domain.GetProfileByAccountID(r.Context(), initiator.GetAccountID())
	if err != nil {
		c.log.WithError(err).Errorf("failed to get profile by user id")
		switch {
		case errors.Is(err, errx.ErrorProfileNotFound):
			c.responser.RenderErr(w, problems.Unauthorized("profile for user does not exist"))
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}

		return
	}

	c.responser.Render(w, http.StatusOK, responses.Profile(res))
}
