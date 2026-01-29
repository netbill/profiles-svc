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

func (s Service) OenProfileUpdateSession(w http.ResponseWriter, r *http.Request) {
	initiator, err := middlewares.AccountData(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	media, profile, err := s.domain.OpenProfileUpdateSession(
		r.Context(),
		initiator.AccountID,
	)
	if err != nil {
		s.log.WithError(err).Errorf("failed to get preload link for update avatar")
		switch {
		case errors.Is(err, errx.ErrorProfileNotFound):
			ape.RenderErr(w, problems.NotFound("profile for user does not exist"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, 200, responses.UpdateProfileSession(media, profile))
}
