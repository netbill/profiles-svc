package controller

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/rest/responses"
	"github.com/netbill/profiles-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationGetProfileByUsername = "get_profile_by_username"

func (c *Controller) GetProfileByUsername(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetProfileByUsername)

	username := chi.URLParam(r, "username")

	res, err := c.modules.Profile.GetByUsername(r.Context(), username)
	switch {
	case errors.Is(err, errx.ErrorProfileNotExists):
		log.Info("profile not found")
		render.ResponseError(w, problems.NotFound("profile for user does not exist"))
	case err != nil:
		log.WithError(err).Error("failed to get profile by username")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.Profile(res))
	}
}
