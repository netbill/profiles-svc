package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/ape"
	"github.com/netbill/ape/problems"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/rest/middlewares"
	"github.com/netbill/profiles-svc/internal/rest/responses"
)

func (s Service) GetPreloadLinkForUpdateAvatar(w http.ResponseWriter, r *http.Request) {
	initiator, err := middlewares.AccountData(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	res, err := s.domain.GetPreloadLinkForUpdateAvatar(
		r.Context(),
		initiator.AccountID,
	)
	if err != nil {
		s.log.WithError(err).Errorf("failed to get preload link for update avatar")
		ape.RenderErr(w, problems.InternalError())

		return
	}

	ape.Render(w, 200, res)
}

func (s Service) AcceptUpdateAvatar(w http.ResponseWriter, r *http.Request) {
	initiator, err := middlewares.AccountData(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	uploadData, err := middlewares.UploadFilesData(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get upload session id")
		ape.RenderErr(w, problems.Unauthorized("failed to get upload session id"))

		return
	}

	res, err := s.domain.AcceptUpdateAvatar(
		r.Context(),
		initiator.AccountID,
		uploadData.SessionID,
	)
	if err != nil {
		s.log.WithError(err).Errorf("failed to accept update avatar")
		ape.RenderErr(w, problems.InternalError())

		return
	}

	ape.Render(w, 200, responses.Profile(res))
}

func (s Service) CancelUpdateAvatar(w http.ResponseWriter, r *http.Request) {
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

	res, err := s.domain.CancelUpdateAvatar(
		r.Context(),
		initiator.AccountID,
		uploadFilesData.SessionID,
	)
	if err != nil {
		s.log.WithError(err).Errorf("failed to cancel update avatar")
		ape.RenderErr(w, problems.InternalError())

		return
	}

	ape.Render(w, 200, responses.Profile(res))
}

func (s Service) DeleteMyProfileAvatar(w http.ResponseWriter, r *http.Request) {
	initiator, err := middlewares.AccountData(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	res, err := s.domain.DeleteProfileAvatar(
		r.Context(),
		initiator.AccountID,
	)
	if err != nil {
		s.log.WithError(err).Errorf("failed to delete profile avatar")
		switch {
		case errors.Is(err, errx.ErrorProfileNotFound):
			ape.RenderErr(w, problems.NotFound("profile not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, 200, responses.Profile(res))
}
