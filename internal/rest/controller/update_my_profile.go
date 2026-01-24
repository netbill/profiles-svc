package controller

import (
	"errors"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/netbill/ape"
	"github.com/netbill/ape/problems"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/core/modules/profile"
	"github.com/netbill/profiles-svc/internal/rest/middlewares"

	"github.com/netbill/profiles-svc/internal/rest/requests"
	"github.com/netbill/profiles-svc/internal/rest/responses"
)

func (s Service) UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	initiator, err := middlewares.AccountData(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.UpdateProfile(r)
	if err != nil {
		s.log.WithError(err).Errorf("invalid create profile request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	if req.Data.Id != initiator.AccountID {
		s.log.WithError(err).Errorf("id in body and initiator id mismatch fir update My profile request")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"id": fmt.Errorf(
				"id in body: %s and initiator id: %s mismatch fir update My profile request",
				req.Data.Id,
				initiator.AccountID,
			),
		})...)
	}

	res, err := s.domain.UpdateProfile(r.Context(), initiator.AccountID, profile.UpdateParams{
		Pseudonym:   req.Data.Attributes.Pseudonym,
		Description: req.Data.Attributes.Description,
	})
	if err != nil {
		s.log.WithError(err).Errorf("failed to update profile")
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
