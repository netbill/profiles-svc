package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/rest/requests"
	"github.com/netbill/profiles-svc/internal/rest/responses"
	"github.com/netbill/profiles-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationUpdateProfileOfficial = "update_profile_official"

func (c *Controller) UpdateProfileOfficial(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationUpdateProfileOfficial)

	req, err := requests.UpdateProfileOfficial(r)
	if err != nil {
		log.WithError(err).Info("invalid update profile official request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	res, err := c.modules.Profile.UpdateOfficial(r.Context(), req.Data.Id, req.Data.Attributes.Official)
	switch {
	case errors.Is(err, errx.ErrorProfileNotExists):
		log.Info("profile for user does not exist")
		c.responser.RenderErr(w, problems.NotFound("profile for user does not exist"))
	case err != nil:
		log.WithError(err).Error("failed to update profile official")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusOK, responses.Profile(res))
	}
}
