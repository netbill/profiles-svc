package controller

import (
	"net/http"

	"github.com/netbill/profiles-svc/internal/rest/contexter"
	"github.com/netbill/restkit/problems"
)

func (c Controller) DeleteUploadProfileAvatar(w http.ResponseWriter, r *http.Request) {
	initiator, err := contexter.AccountData(r.Context())
	if err != nil {
		c.log.WithError(err).Error("failed to get user from context")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	uploadFilesData, err := contexter.UploadContentData(r.Context())
	if err != nil {
		c.log.WithError(err).Error("failed to get upload session id")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get upload session id"))

		return
	}

	err = c.domain.DeleteUploadProfileAvatarInSession(
		r.Context(),
		initiator.GetAccountID(),
		uploadFilesData.GetUploadSessionID(),
	)
	if err != nil {
		c.log.WithError(err).Errorf("failed to cancel update avatar")
		c.responser.RenderErr(w, problems.InternalError())

		return
	}

	c.responser.Render(w, 200, nil)
}
