package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/rest/responses"
	"github.com/netbill/profiles-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationOpenUpdateProfileSession = "open_update_profile_session"

func (c *Controller) OenUpdateProfileSession(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationOpenUpdateProfileSession)

	media, profile, err := c.modules.Profile.OpenUpdateSession(
		r.Context(),
		scope.AccountActor(r),
	)

	switch {
	case errors.Is(err, errx.ErrorProfileNotExists):
		log.Info("profile for user does not exist")
		c.responser.RenderErr(w, problems.Unauthorized("profile for user does not exist"))
	case err != nil:
		log.WithError(err).Error("failed to open update profile session")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusOK, responses.UpdateProfileSession(media, profile))
	}
}
