package rest

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/netbill/logium"
	"github.com/netbill/profiles-svc/internal"
	"github.com/netbill/restkit/auth/roles"
)

type Handlers interface {
	GetMyProfile(w http.ResponseWriter, r *http.Request)

	GetProfileByUsername(w http.ResponseWriter, r *http.Request)
	GetProfileByID(w http.ResponseWriter, r *http.Request)

	FilterProfiles(w http.ResponseWriter, r *http.Request)

	UpdateMyProfile(w http.ResponseWriter, r *http.Request)

	UpdateProfileOfficial(w http.ResponseWriter, r *http.Request)
	UpdateProfileUsername(w http.ResponseWriter, r *http.Request)
}
type Middlewares interface {
	Auth() func(http.Handler) http.Handler
	RoleGrant(allowedRoles map[string]bool) func(http.Handler) http.Handler
}

type Service struct {
	handlers    Handlers
	middlewares Middlewares
	log         logium.Logger
}

func New(
	log logium.Logger,
	middlewares Middlewares,
	handlers Handlers,
) *Service {
	return &Service{
		log:         log,
		middlewares: middlewares,
		handlers:    handlers,
	}
}

func (s *Service) Run(ctx context.Context, cfg internal.Config) {
	auth := s.middlewares.Auth()
	sysmoder := s.middlewares.RoleGrant(map[string]bool{
		roles.SystemAdmin: true,
		roles.SystemModer: true,
	})

	r := chi.NewRouter()

	r.Route("/profiles-svc", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/profiles", func(r chi.Router) {
				r.Get("/", s.handlers.FilterProfiles)

				r.Get("/u/{username}", s.handlers.GetProfileByUsername)

				r.With(auth).Route("/me", func(r chi.Router) {
					r.Get("/", s.handlers.GetMyProfile)
					r.Put("/", s.handlers.UpdateMyProfile)
					r.Patch("/username", s.handlers.UpdateProfileUsername)
				})

				r.Route("/{user_id}", func(r chi.Router) {
					r.Get("/", s.handlers.GetProfileByID)

					r.With(auth, sysmoder).Patch("/official", s.handlers.UpdateProfileOfficial)
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

	s.log.Infof("starting REST service on %s", cfg.Rest.Port)

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
		s.log.Warnf("shutting down REST service...")
	case err := <-errCh:
		if err != nil {
			s.log.Errorf("REST server error: %v", err)
		}
	}

	shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shCtx); err != nil {
		s.log.Errorf("REST shutdown error: %v", err)
	} else {
		s.log.Warnf("REST server stopped")
	}
}
