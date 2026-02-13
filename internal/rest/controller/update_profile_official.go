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

func (c *Controller) UpdateProfileOfficial(w http.ResponseWriter, r *http.Request) {
	req, err := requests.UpdateProfileOfficial(r)
	if err != nil {
		scope.Log(r).WithError(err).Errorf("invalid update official request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	res, err := c.modules.Profile.UpdateOfficial(r.Context(), req.Data.Id, req.Data.Attributes.Official)
	switch {
	case errors.Is(err, errx.ErrorProfileNotExists):
		scope.Log(r).Info("profile for user does not exist")
		c.responser.RenderErr(w, problems.NotFound("profile for user does not exist"))
	case err != nil:
		scope.Log(r).WithError(err).Error("failed to update profile official")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusOK, responses.Profile(res))
	}
}
