package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/rest/responses"
	"github.com/netbill/profiles-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	res, err := c.modules.Profile.GetByAccountID(r.Context(), scope.AccountAuthClaims(r).GetAccountID())
	switch {
	case errors.Is(err, errx.ErrorProfileNotExists):
		scope.Log(r).Infof("profile for user does not exist")
		c.responser.RenderErr(w, problems.Unauthorized("profile for user does not exist"))
	case err != nil:
		scope.Log(r).Errorf("failed to get profile by user id")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusOK, responses.Profile(res))
	}
}
