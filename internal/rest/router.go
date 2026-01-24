package rest

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/netbill/logium"
	"github.com/netbill/restkit/tokens/roles"
)

type Handlers interface {
	GetMyProfile(w http.ResponseWriter, r *http.Request)

	GetProfileByUsername(w http.ResponseWriter, r *http.Request)
	GetProfileByID(w http.ResponseWriter, r *http.Request)

	FilterProfiles(w http.ResponseWriter, r *http.Request)

	UpdateMyProfile(w http.ResponseWriter, r *http.Request)
	UpdateProfileOfficial(w http.ResponseWriter, r *http.Request)

	GetPreloadLinkForUpdateAvatar(w http.ResponseWriter, r *http.Request)
	AcceptUpdateAvatar(w http.ResponseWriter, r *http.Request)
	CancelUpdateAvatar(w http.ResponseWriter, r *http.Request)
	DeleteMyProfileAvatar(w http.ResponseWriter, r *http.Request)
}
type Middlewares interface {
	AccountAuth() func(http.Handler) http.Handler
	AccountRolesGrant(allowedRoles map[string]bool) func(http.Handler) http.Handler
	ConfirmUploadFiles(scope string) func(http.Handler) http.Handler
}

type Service struct {
	handlers    Handlers
	middlewares Middlewares
	log         *logium.Logger
}

func New(
	log *logium.Logger,
	middlewares Middlewares,
	handlers Handlers,
) *Service {
	return &Service{
		log:         log,
		middlewares: middlewares,
		handlers:    handlers,
	}
}

type Config struct {
	Port              string
	TimeoutRead       time.Duration
	TimeoutReadHeader time.Duration
	TimeoutWrite      time.Duration
	TimeoutIdle       time.Duration
}

func (s *Service) Run(ctx context.Context, cfg Config) {
	auth := s.middlewares.AccountAuth()
	sysmoder := s.middlewares.AccountRolesGrant(map[string]bool{
		roles.SystemAdmin: true,
		roles.SystemModer: true,
	})
	uploadProfileAvatar := s.middlewares.ConfirmUploadFiles("upload_profile_avatar")

	r := chi.NewRouter()

	// CORS for swagger UI documentation need to delete after configuring nginx
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5002"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/profiles-svc", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/profiles", func(r chi.Router) {
				r.Get("/", s.handlers.FilterProfiles)

				r.Get("/u/{username}", s.handlers.GetProfileByUsername)

				r.With(auth).Route("/me", func(r chi.Router) {
					r.Get("/", s.handlers.GetMyProfile)
					r.Put("/", s.handlers.UpdateMyProfile)
					r.Route("/avatar", func(r chi.Router) {
						r.Route("/upload", func(r chi.Router) {
							r.Get("/", s.handlers.GetPreloadLinkForUpdateAvatar)

							r.With(uploadProfileAvatar).Post("/", s.handlers.AcceptUpdateAvatar)
							r.With(uploadProfileAvatar).Delete("/", s.handlers.CancelUpdateAvatar)
						})

						r.Delete("/", s.handlers.DeleteMyProfileAvatar)
					})
				})

				r.Route("/{account_id}", func(r chi.Router) {
					r.Get("/", s.handlers.GetProfileByID)

					r.With(auth, sysmoder).Patch("/official", s.handlers.UpdateProfileOfficial)
				})
			})
		})
	})

	srv := &http.Server{
		Addr:              cfg.Port,
		Handler:           r,
		ReadTimeout:       cfg.TimeoutRead,
		ReadHeaderTimeout: cfg.TimeoutReadHeader,
		WriteTimeout:      cfg.TimeoutWrite,
		IdleTimeout:       cfg.TimeoutIdle,
	}

	s.log.Infof("starting REST service on %s", cfg.Port)

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
