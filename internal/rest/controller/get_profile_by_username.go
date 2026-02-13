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

const operationGetProfileByUsername = "get_profile_by_username"

func (c *Controller) GetProfileByUsername(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetProfileByUsername)

	username := chi.URLParam(r, "username")

	res, err := c.modules.Profile.GetByUsername(r.Context(), username)
	switch {
	case errors.Is(err, errx.ErrorProfileNotExists):
		log.Info("profile not found")
		c.responser.RenderErr(w, problems.NotFound("profile for user does not exist"))
	case err != nil:
		log.WithError(err).Error("failed to get profile by username")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusOK, responses.Profile(res))
	}
}
