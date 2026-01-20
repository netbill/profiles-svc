package controller

import (
	"errors"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/rest/requests"
	"github.com/netbill/profiles-svc/internal/rest/responses"
	"github.com/netbill/restkit/ape"
	"github.com/netbill/restkit/ape/problems"
)

func (s Service) UpdateProfileUsername(w http.ResponseWriter, r *http.Request) {
	req, err := requests.UpdateProfileUsername(r)
	if err != nil {
		s.log.WithError(err).Errorf("invalid update username request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	res, err := s.domain.UpdateProfileUsername(r.Context(), req.Data.Id, req.Data.Attributes.Username)
	if err != nil {
		s.log.WithError(err).Errorf("failed to update username status")
		switch {
		case errors.Is(err, errx.ErrorProfileNotFound):
			ape.RenderErr(w, problems.NotFound("profile for user does not exist"))
		case errors.Is(err, errx.ErrorUsernameAlreadyTaken):
			ape.RenderErr(w, problems.Conflict("username is already taken"))
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
