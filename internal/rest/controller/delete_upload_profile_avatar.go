package controller

import (
	"net/http"

	"github.com/netbill/ape"
	"github.com/netbill/ape/problems"
	"github.com/netbill/profiles-svc/internal/rest/middlewares"
)

func (s Service) DeleteUploadProfileAvatar(w http.ResponseWriter, r *http.Request) {
	initiator, err := middlewares.AccountData(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	uploadFilesData, err := middlewares.UploadFilesData(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get upload session id")
		ape.RenderErr(w, problems.Unauthorized("failed to get upload session id"))

		return
	}

	err = s.domain.DeleteUploadProfileAvatarInSession(
		r.Context(),
		initiator.AccountID,
		uploadFilesData.UploadSessionID,
	)
	if err != nil {
		s.log.WithError(err).Errorf("failed to cancel update avatar")
		ape.RenderErr(w, problems.InternalError())

		return
	}

	ape.Render(w, 200, nil)
}
