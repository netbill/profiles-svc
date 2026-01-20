package controller

import (
	"errors"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/core/modules/profile"
	"github.com/netbill/profiles-svc/internal/rest"
	"github.com/netbill/profiles-svc/internal/rest/requests"
	"github.com/netbill/profiles-svc/internal/rest/responses"
	"github.com/netbill/restkit/ape"
	"github.com/netbill/restkit/ape/problems"
)

func (s Service) CreateMyProfile(w http.ResponseWriter, r *http.Request) {
	initiator, err := rest.AccountData(r.Context())
	if err != nil {
		s.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.CreateProfile(r)
	if err != nil {
		s.log.WithError(err).Errorf("invalid create profile request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	if req.Data.Id != initiator.ID {
		s.log.WithError(err).Errorf("id in body and initiator id mismatch fir create My profile request")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"id": fmt.Errorf(
				"id in body: %s and initiator id: %s mismatch fir create My profile request",
				req.Data.Id,
				initiator.ID,
			),
		})...)
	}

	res, err := s.domain.CreateProfile(r.Context(), initiator.ID, profile.CreateParams{
		Pseudonym:   req.Data.Attributes.Pseudonym,
		Description: req.Data.Attributes.Description,
		Avatar:      req.Data.Attributes.Avatar,
	})
	if err != nil {
		s.log.WithError(err).Errorf("failed to create profile")
		switch {
		case errors.Is(err, errx.ErrorProfileNotFound):
			ape.RenderErr(w, problems.Unauthorized("profile for user does not exist"))
		case errors.Is(err, errx.ErrorUsernameAlreadyTaken):
			ape.RenderErr(w, problems.Conflict("user name is already taken"))
		case errors.Is(err, errx.ErrorUsernameIsNotAllowed):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"username": err,
			})...)

		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Profile(res))
}
