package controller

import (
	"net/http"

	"github.com/netbill/profiles-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationDeleteUploadProfileAvatar = "delete_upload_profile_avatar"

func (c *Controller) DeleteUploadProfileAvatar(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeleteUploadProfileAvatar)

	err := c.modules.Profile.DeleteUploadAvatar(
		r.Context(),
		scope.AccountActor(r),
		scope.UploadScope(r),
	)
	switch {
	case err != nil:
		log.WithError(err).Error("failed to delete upload avatar")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		log.Info("upload avatar deleted")
		c.responser.Render(w, http.StatusOK, nil)
	}
}
