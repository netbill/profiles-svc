package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/rest/responses"
	"github.com/netbill/profiles-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) OenUpdateProfileSession(w http.ResponseWriter, r *http.Request) {
	media, profile, err := c.modules.Profile.OpenUpdateSession(r.Context(), scope.AccountAuthClaims(r).GetAccountID())
	switch {
	case errors.Is(err, errx.ErrorProfileNotExists):
		scope.Log(r).Errorf("profile for user does not exist")
		c.responser.RenderErr(w, problems.Unauthorized("profile for user does not exist"))
	case err != nil:
		scope.Log(r).Errorf("failed to get preload link for update avatar")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, 200, responses.UpdateProfileSession(media, profile))
	}
}
