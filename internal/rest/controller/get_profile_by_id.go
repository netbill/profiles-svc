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

func (c *Controller) GetProfileByID(w http.ResponseWriter, r *http.Request) {
	accountID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		scope.Log(r).WithError(err).Errorf("invalid account id")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid account id: %s", chi.URLParam(r, "account_id")),
		})...)

		return
	}

	res, err := c.modules.Profile.GetByAccountID(r.Context(), accountID)
	switch {
	case errors.Is(err, errx.ErrorProfileNotExists):
		scope.Log(r).Infof("profile for user does not exist")
		c.responser.RenderErr(w, problems.NotFound("profile for user does not exist"))
	case err != nil:
		scope.Log(r).Errorf("failed to get profile by user id")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusOK, responses.Profile(res))
	}
}
