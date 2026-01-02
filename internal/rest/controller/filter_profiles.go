package controller

import (
	"net/http"
	"strings"

	"github.com/netbill/ape"
	"github.com/netbill/ape/problems"
	"github.com/netbill/pagi"
	"github.com/netbill/profiles-svc/internal/domain/modules/profile"
	"github.com/netbill/profiles-svc/internal/rest/responses"
)

func (s Service) FilterProfiles(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, offset := pagi.GetPagination(r)

	filters := profile.FilterParams{}

	if usernameLike := strings.TrimSpace(q.Get("username_like")); usernameLike != "" {
		filters.UsernamePrefix = &usernameLike
	}

	if pseudonym := strings.TrimSpace(q.Get("pseudonym")); pseudonym != "" {
		filters.PseudonymPrefix = &pseudonym
	}

	res, err := s.domain.FilterProfile(r.Context(), filters, limit, offset)
	if err != nil {
		s.log.WithError(err).Error("failed to filter profiles")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, http.StatusOK, responses.ProfileCollection(res))
}
