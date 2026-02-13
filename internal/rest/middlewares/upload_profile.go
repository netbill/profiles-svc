package middlewares

import (
	"net/http"

	"github.com/netbill/profiles-svc/internal/rest/scope"
	"github.com/netbill/restkit/headers"
	"github.com/netbill/restkit/problems"
)

func (p *Provider) UpdateOwnProfileMediaContent() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			initiator := scope.AccountAuthClaims(r)

			token, err := headers.GetUploadContent(r)
			if err != nil {
				scope.Log(r).WithAccountAuthClaims(initiator).Debug("upload token missing")
				p.responser.RenderErr(w, problems.Unauthorized("failed to get token"))

				return
			}

			uploadClaims, err := p.tokenManager.ParseUploadProfileContentToken(token)
			if err != nil {
				scope.Log(r).WithAccountAuthClaims(initiator).Info("upload token invalid")
				p.responser.RenderErr(w, problems.Unauthorized("invalid upload profile token"))

				return
			}

			if uploadClaims.GetAccountID() != initiator.GetAccountID() {
				scope.Log(r).Info("account is not owner of the profile")
				p.responser.RenderErr(w, problems.Unauthorized("account is not owner of the profile"))

				return
			}

			if uploadClaims.GetResourceID() != initiator.GetAccountID().String() {
				scope.Log(r).Info("upload token is not for profile content")
				p.responser.RenderErr(w, problems.Unauthorized("invalid upload profile token"))

				return
			}

			next.ServeHTTP(w, r.WithContext(scope.CtxUploadContent(r.Context(), uploadClaims)))
		})
	}
}
