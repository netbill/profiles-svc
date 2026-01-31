package controller

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/rest/responses"
	"github.com/netbill/restkit/problems"
)

func (c Controller) GetProfileByUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	res, err := c.domain.GetProfileByUsername(r.Context(), username)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get profile by username")
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
