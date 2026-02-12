package controller

import (
	"net/http"

	"github.com/netbill/profiles-svc/internal/rest/contexter"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) CancelUpdateProfileSession(w http.ResponseWriter, r *http.Request) {
	initiator, err := contexter.AccountData(r.Context())
	if err != nil {
		c.Log(r).WithError(err).Error("failed to get user from context")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	uploadFilesData, err := contexter.UploadProfileContentData(r.Context())
	if err != nil {
		c.Log(r).WithError(err).Error("failed to get upload session id")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get upload session id"))

		return
	}

	err = c.modules.Profile.CancelUpdateSession(
		r.Context(),
		initiator.GetAccountID(),
		uploadFilesData.GetSessionID(),
	)
	switch {
	case err != nil:
		c.Log(r).WithError(err).Errorf("failed to cancel update profile session")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.Log(r).Info("profile update session cancelled successfully")
		c.responser.Render(w, http.StatusOK, nil)
	}
}
