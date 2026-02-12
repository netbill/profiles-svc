package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/rest/contexter"
	"github.com/netbill/profiles-svc/internal/rest/responses"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	initiator, err := contexter.AccountData(r.Context())
	if err != nil {
		c.Log(r).WithError(err).Error("failed to get account from context")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get account from context"))

		return
	}

	res, err := c.modules.Profile.GetByAccountID(r.Context(), initiator.GetAccountID())
	switch {
	case errors.Is(err, errx.ErrorProfileNotExists):
		c.Log(r).Infof("profile for user does not exist")
		c.responser.RenderErr(w, problems.Unauthorized("profile for user does not exist"))
	case err != nil:
		c.Log(r).Errorf("failed to get profile by user id")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusOK, responses.Profile(res))
	}
}
