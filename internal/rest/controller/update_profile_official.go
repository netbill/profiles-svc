package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/ape"
	"github.com/netbill/ape/problems"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/rest/requests"
	"github.com/netbill/profiles-svc/internal/rest/responses"
)

func (s Service) UpdateProfileOfficial(w http.ResponseWriter, r *http.Request) {
	req, err := requests.UpdateProfileOfficial(r)
	if err != nil {
		s.log.WithError(err).Errorf("invalid update official request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	res, err := s.domain.UpdateProfileOfficial(r.Context(), req.Data.Id, req.Data.Attributes.Official)
	if err != nil {
		s.log.WithError(err).Errorf("failed to update official status")
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
