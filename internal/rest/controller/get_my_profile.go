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

func (s Service) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	initiator, err := middlewares.AccountData(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get account from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get account from context"))

		return
	}

	res, err := s.domain.GetProfileByID(r.Context(), initiator.AccountID)
	if err != nil {
		s.log.WithError(err).Errorf("failed to get profile by user id")
		switch {
		case errors.Is(err, errx.ErrorProfileNotFound):
			ape.RenderErr(w, problems.Unauthorized("profile for user does not exist"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Profile(res))
}
