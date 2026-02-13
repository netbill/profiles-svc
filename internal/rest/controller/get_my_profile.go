package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/rest/responses"
	"github.com/netbill/profiles-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationGetMyProfile = "get_my_profile"

func (c *Controller) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetMyProfile)

	res, err := c.modules.Profile.GetByAccountID(r.Context(), scope.AccountActor(r))
	switch {
	case errors.Is(err, errx.ErrorProfileNotExists):
		log.Info("profile for user does not exist")
		c.responser.RenderErr(w, problems.Unauthorized("profile for user does not exist"))
	case err != nil:
		log.WithError(err).Error("failed to get profile by account id")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusOK, responses.Profile(res))
	}
}
