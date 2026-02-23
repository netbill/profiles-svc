package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/netbill/profiles-svc/pkg/log"
	"github.com/netbill/restkit/tokens"
)

type Handlers interface {
	GetMyProfile(w http.ResponseWriter, r *http.Request)

	GetProfileByUsername(w http.ResponseWriter, r *http.Request)
	GetProfileByID(w http.ResponseWriter, r *http.Request)

	FilterProfiles(w http.ResponseWriter, r *http.Request)

	UpdateMyProfile(w http.ResponseWriter, r *http.Request)
	UpdateProfileOfficial(w http.ResponseWriter, r *http.Request)

	CreateMyProfileUploadMediaLink(w http.ResponseWriter, r *http.Request)
	DeleteMyProfileUploadAvatar(w http.ResponseWriter, r *http.Request)
}

type Middlewares interface {
	AccountAuth(
		allowedRoles ...string,
	) func(next http.Handler) http.Handler
	Logger(log *log.Logger) func(next http.Handler) http.Handler
	CorsDocs() func(next http.Handler) http.Handler
}

type Server struct {
	handlers    Handlers
	middlewares Middlewares
}

func New(m Middlewares, h Handlers) *Server {
	return &Server{
		middlewares: m,
		handlers:    h,
	}
}

type Config struct {
	Port              int
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

func (s *Server) Run(ctx context.Context, log *log.Logger, cfg Config) {
	auth := s.middlewares.AccountAuth()
	sysmoder := s.middlewares.AccountAuth(tokens.RoleSystemAdmin, tokens.RoleSystemModer)

	r := chi.NewRouter()
	r.Use(
		s.middlewares.Logger(log),
		s.middlewares.CorsDocs(),
	)

	r.Route("/profiles-svc", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/profiles", func(r chi.Router) {
				r.Get("/", s.handlers.FilterProfiles)

				r.With(auth).Route("/me", func(r chi.Router) {
					r.Get("/", s.handlers.GetMyProfile)
					r.Put("/", s.handlers.UpdateMyProfile)

					r.Route("/media", func(r chi.Router) {
						r.Route("/upload", func(r chi.Router) {
							r.Post("/url", s.handlers.CreateMyProfileUploadMediaLink)

							r.Delete("/avatar", s.handlers.DeleteMyProfileUploadAvatar)
						})
					})
				})

				r.Get("/@{username}", s.handlers.GetProfileByUsername)

				r.Route("/{account_id:[0-9a-fA-F-]{36}}", func(r chi.Router) {
					r.Get("/", s.handlers.GetProfileByID)
					r.With(sysmoder).Patch("/official", s.handlers.UpdateProfileOfficial)
				})
			})
		})
	})

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           r,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}

	log.WithField("port", cfg.Port).Info("starting http service...")

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
		log.Info("shutting down http service...")
	case err := <-errCh:
		if err != nil {
			log.WithError(err).Error("http server error")
		}
	}

	shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shCtx); err != nil {
		log.WithError(err).Error("failed to shutdown http server gracefully")
	} else {
		log.Info("http server shutdown gracefully")
	}
}
