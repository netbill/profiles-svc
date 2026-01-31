package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/rest/requests"
	"github.com/netbill/profiles-svc/internal/rest/responses"
	"github.com/netbill/restkit/problems"
)

func (c Controller) UpdateProfileOfficial(w http.ResponseWriter, r *http.Request) {
	req, err := requests.UpdateProfileOfficial(r)
	if err != nil {
		c.log.WithError(err).Errorf("invalid update official request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	res, err := c.domain.UpdateProfileOfficial(r.Context(), req.Data.Id, req.Data.Attributes.Official)
	if err != nil {
		c.log.WithError(err).Errorf("failed to update official status")
		switch {
		case errors.Is(err, errx.ErrorProfileNotFound):
			c.responser.RenderErr(w, problems.NotFound("profile for user does not exist"))
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}

		return
	}

	c.responser.Render(w, http.StatusOK, responses.Profile(res))
}
