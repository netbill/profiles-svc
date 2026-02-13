package controller

import (
	"net/http"

	"github.com/netbill/profiles-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) CancelUpdateProfileSession(w http.ResponseWriter, r *http.Request) {
	err := c.modules.Profile.CancelUpdateSession(
		r.Context(),
		scope.AccountAuthClaims(r).GetAccountID(),
		scope.UploadContentClaims(r).GetSessionID(),
	)

	switch {
	case err != nil:
		scope.Log(r).WithError(err).Errorf("failed to cancel update profile session")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		scope.Log(r).Info("profile update session cancelled successfully")
		c.responser.Render(w, http.StatusOK, nil)
	}
}
