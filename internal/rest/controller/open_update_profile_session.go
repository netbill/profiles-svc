package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/rest/contexter"
	"github.com/netbill/profiles-svc/internal/rest/responses"
	"github.com/netbill/restkit/problems"
)

func (c Controller) OenProfileUpdateSession(w http.ResponseWriter, r *http.Request) {
	initiator, err := contexter.AccountData(r.Context())
	if err != nil {
		c.log.WithError(err).Error("failed to get user from context")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	media, profile, err := c.domain.OpenProfileUpdateSession(
		r.Context(),
		initiator.GetAccountID(),
	)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get preload link for update avatar")
		switch {
		case errors.Is(err, errx.ErrorProfileNotFound):
			c.responser.RenderErr(w, problems.NotFound("profile for user does not exist"))
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}

		return
	}

	c.responser.Render(w, 200, responses.UpdateProfileSession(media, profile))
}
