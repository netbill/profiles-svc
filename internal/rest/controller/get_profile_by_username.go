package controller

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/rest/responses"
	"github.com/netbill/profiles-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) GetProfileByUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	res, err := c.modules.Profile.GetByUsername(r.Context(), username)
	switch {
	case errors.Is(err, errx.ErrorProfileNotExists):
		scope.Log(r).Infof("profile for user does not exist")
		c.responser.RenderErr(w, problems.NotFound("profile for user does not exist"))
	case err != nil:
		scope.Log(r).Errorf("failed to get profile by username")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusOK, responses.Profile(res))
	}
}
