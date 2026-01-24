package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/ape"
	"github.com/netbill/ape/problems"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/rest/responses"
)

func (s Service) GetProfileByID(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		s.log.WithError(err).Errorf("invalid account id")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid account id: %s", chi.URLParam(r, "account_id")),
		})...)

		return
	}

	res, err := s.domain.GetProfileByID(r.Context(), userID)
	if err != nil {
		s.log.WithError(err).Errorf("failed to get profile by user id")
		switch {
		case errors.Is(err, errx.ErrorProfileNotFound):
			ape.RenderErr(w, problems.NotFound("profile for user does not exist"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Profile(res))
}
