package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/rest/responses"
	"github.com/netbill/profiles-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationGetProfileByID = "get_profile_by_id"

func (c *Controller) GetProfileByID(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetProfileByID)

	accountID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		log.WithError(err).Warn("invalid account id")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"path": fmt.Errorf("invalid account id: %s", chi.URLParam(r, "account_id")),
		})...)
		return
	}

	res, err := c.modules.Profile.GetMy(r.Context(), accountID)
	switch {
	case errors.Is(err, errx.ErrorProfileNotExists):
		log.Info("profile not found")
		c.responser.RenderErr(w, problems.NotFound("profile for user does not exist"))
	case err != nil:
		log.WithError(err).Error("failed to get profile by account id")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusOK, responses.Profile(res))
	}
}
