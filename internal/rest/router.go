package rest

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/netbill/logium"
	"github.com/netbill/profiles-svc/internal"
	"github.com/netbill/profiles-svc/internal/rest/meta"
	"github.com/netbill/restkit/roles"
)

type Handlers interface {
	GetMyProfile(w http.ResponseWriter, r *http.Request)

	//CreateMyProfile(w http.ResponseWriter, r *http.Request)
	GetProfileByUsername(w http.ResponseWriter, r *http.Request)
	GetProfileByID(w http.ResponseWriter, r *http.Request)

	FilterProfiles(w http.ResponseWriter, r *http.Request)

	UpdateMyProfile(w http.ResponseWriter, r *http.Request)
	//UpdateMyUsername(w http.ResponseWriter, r *http.Request)
	UpdateOfficial(w http.ResponseWriter, r *http.Request)

	//ResetProfile(w http.ResponseWriter, r *http.Request)
}

type Middleware interface {
	Auth(userCtxKey interface{}, skUser string) func(http.Handler) http.Handler
	RoleGrant(userCtxKey interface{}, allowedRoles map[string]bool) func(http.Handler) http.Handler
}

func Run(ctx context.Context, cfg internal.Config, log logium.Logger, m Middleware, h Handlers) {
	auth := m.Auth(meta.AccountDataCtxKey, cfg.JWT.User.AccessToken.SecretKey)
	sysmoder := m.RoleGrant(meta.AccountDataCtxKey, map[string]bool{
		roles.SystemModer: true,
		roles.SystemAdmin: true,
	})

	r := chi.NewRouter()

	r.Route("/profiles-svc", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/profiles", func(r chi.Router) {
				r.Get("/", h.FilterProfiles)
				r.Get("/u/{username}", h.GetProfileByUsername)

				r.With(auth).Route("/me", func(r chi.Router) {
					r.Get("/", h.GetMyProfile)
					r.Put("/", h.UpdateMyProfile)
				})

				r.Route("/{user_id}", func(r chi.Router) {
					r.Get("/", h.GetProfileByID)

					r.With(auth, sysmoder).Patch("/official", h.UpdateOfficial)
					//r.With(auth, sysmoder).Put("/reset", h.ResetProfile)
				})
			})
		})
	})

	srv := &http.Server{
		Addr:              cfg.Rest.Port,
		Handler:           r,
		ReadTimeout:       cfg.Rest.Timeouts.Read,
		ReadHeaderTimeout: cfg.Rest.Timeouts.ReadHeader,
		WriteTimeout:      cfg.Rest.Timeouts.Write,
		IdleTimeout:       cfg.Rest.Timeouts.Idle,
	}

	log.Infof("starting REST service on %s", cfg.Rest.Port)

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		} else {
			errCh <- nil
		}
	}()

	select {
	case <-ctx.Done():
		log.Info("shutting down REST service...")
	case err := <-errCh:
		if err != nil {
			log.Errorf("REST server error: %v", err)
		}
	}

	shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shCtx); err != nil {
		log.Errorf("REST shutdown error: %v", err)
	} else {
		log.Info("REST server stopped")
	}
}
