package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/rest/requests"
	"github.com/netbill/profiles-svc/internal/rest/responses"
	"github.com/netbill/profiles-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationUpdateProfileOfficial = "update_profile_official"

func (c *Controller) UpdateProfileOfficial(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationUpdateProfileOfficial)

	req, err := requests.UpdateProfileOfficial(r)
	if err != nil {
		log.WithError(err).Info("invalid update profile official request")
		render.ResponseError(w, problems.BadRequest(err)...)
		return
	}

	res, err := c.modules.Profile.UpdateOfficial(r.Context(), req.Data.Id, req.Data.Attributes.Official)
	switch {
	case errors.Is(err, errx.ErrorProfileNotExists):
		log.Info("profile for user does not exist")
		render.ResponseError(w, problems.NotFound("profile for user does not exist"))
	case err != nil:
		log.WithError(err).Error("failed to update profile official")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.Profile(res))
	}
}
