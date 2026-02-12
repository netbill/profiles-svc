package middlewares

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/netbill/logium"
	"github.com/netbill/profiles-svc/internal/rest/contexter"
	"github.com/netbill/restkit/headers"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/tokens"
)

type responser interface {
	Render(w http.ResponseWriter, status int, res ...interface{})
	RenderErr(w http.ResponseWriter, errs ...error)
}

type tokenManager interface {
	ParseAccessClaims(token string) (tokens.AccountClaims, error)

	ParseUploadProfileContentToken(token string) (tokens.UploadContentClaims, error)
	GenerateUploadProfileMediaToken(
		OwnerAccountID uuid.UUID,
		UploadSessionID uuid.UUID,
	) (string, error)
}

type Provider struct {
	log *logium.Logger

	tokenManager tokenManager
	responser    responser
}

func New(
	log *logium.Logger,
	responser responser,
	tokenManager tokenManager,
) *Provider {
	return &Provider{
		log:          log,
		tokenManager: tokenManager,
		responser:    responser,
	}
}

func (p *Provider) AccountAuth(allowedRoles ...string) func(next http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(allowedRoles))
	for _, r := range allowedRoles {
		allowed[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := headers.GetAuthorizationToken(r)
			if err != nil {
				p.log.WithError(err).WithRequest(r).Debug("account authentication failed")
				p.responser.RenderErr(w, problems.Unauthorized("account authentication failed"))

				return
			}

			claims, err := p.tokenManager.ParseAccessClaims(token)
			if err != nil {
				p.log.WithError(err).WithRequest(r).Info("account authentication failed")
				p.responser.RenderErr(w, problems.Unauthorized("account authentication failed"))

				return
			}

			if len(allowed) > 0 {
				if _, ok := allowed[claims.Role]; !ok {
					p.log.WithRequest(r).WithAccount(claims).Debug("account authentication rejected by role")
					p.responser.RenderErr(w, problems.Unauthorized("invalid authentication role"))

					return
				}
			}

			ctx := context.WithValue(r.Context(), contexter.AccountDataCtxKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (p *Provider) UpdateOwnProfile() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			initiator, err := contexter.AccountData(r.Context())
			if err != nil {
				p.log.WithError(err).WithRequest(r).Warn("authorization context is missing")
				p.responser.RenderErr(w, problems.Unauthorized("failed to get user from context"))

				return
			}

			token, err := headers.GetUploadToken(r)
			if err != nil {
				p.log.WithError(err).WithRequest(r).WithAccount(initiator).Debug("upload token missing")
				p.responser.RenderErr(w, problems.Unauthorized("failed to get token"))

				return
			}

			uploadClaims, err := p.tokenManager.ParseUploadProfileContentToken(token)
			if err != nil {
				p.log.WithError(err).WithRequest(r).WithAccount(initiator).Info("upload token invalid")
				p.responser.RenderErr(w, problems.Unauthorized("invalid upload profile token"))

				return
			}

			if initiator.GetAccountID().String() != uploadClaims.ResourceID {
				p.log.WithRequest(r).WithAccount(initiator).WithUploadSession(uploadClaims).
					Warn("upload token owner it's not have an account with this profile")

				p.responser.RenderErr(w, problems.Unauthorized("invalid upload profile token"))
				return
			}

			ctx := context.WithValue(r.Context(), contexter.UploadContentCtxKey, uploadClaims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
