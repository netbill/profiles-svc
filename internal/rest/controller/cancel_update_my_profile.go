package controller

import (
	"net/http"

	"github.com/netbill/profiles-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationCancelUpdateMyProfile = "cancel_update_my_profile"

func (c *Controller) CancelUpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationCancelUpdateMyProfile)

	err := c.modules.Profile.CancelUpdateSession(
		r.Context(),
		scope.AccountActor(r),
		scope.UploadScope(r),
	)
	switch {
	case err != nil:
		log.WithError(err).Error("failed to cancel update profile session")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		log.Info("profile update session cancelled")
		c.responser.Render(w, http.StatusOK, nil)
	}
}
