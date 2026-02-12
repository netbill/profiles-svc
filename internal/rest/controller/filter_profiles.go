package controller

import (
	"net/http"
	"strings"

	"github.com/netbill/profiles-svc/internal/core/modules/profile"
	"github.com/netbill/profiles-svc/internal/rest/responses"
	"github.com/netbill/restkit/pagi"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) FilterProfiles(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, offset := pagi.GetPagination(r)

	filters := profile.FilterParams{}

	if text := strings.TrimSpace(q.Get("text")); text != "" {
		filters.Text = &text
	}
	if official := q.Get("official"); official != "" {
		officialBool := official == "true"
		filters.Official = &officialBool
	}

	res, err := c.modules.Profile.GetList(r.Context(), filters, limit, offset)
	switch {
	case err != nil:
		c.Log(r).WithError(err).Error("failed to filter profiles")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusOK, responses.ProfileCollection(r, res))
	}
}
