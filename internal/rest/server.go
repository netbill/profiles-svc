package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/netbill/profiles-svc/internal/media"
	"github.com/netbill/profiles-svc/pkg/log"
)

type profileController interface {
	GetMy(w http.ResponseWriter, r *http.Request)

	GetByUsername(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)

	Filter(w http.ResponseWriter, r *http.Request)

	UpdateMy(w http.ResponseWriter, r *http.Request)

	CreateUploadMediaLink(w http.ResponseWriter, r *http.Request)
	DeleteUploadMedia(w http.ResponseWriter, r *http.Request)
}

type Middlewares interface {
	AccountAuth(
		allowedRoles ...string,
	) func(next http.Handler) http.Handler
	Logger(log *log.Logger) func(next http.Handler) http.Handler
	CorsDocs() func(next http.Handler) http.Handler
	ResolverUrl(resolver *media.Resolver) func(next http.Handler) http.Handler
}

type Server struct {
	middlewares Middlewares

	profile profileController

	log      *log.Logger
	resolver *media.Resolver
	config   Config
}

type Config struct {
	Port              int
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

type ServerDeps struct {
	Middlewares Middlewares

	Profile profileController

	Log      *log.Logger
	Resolver *media.Resolver
}

func NewServer(deps ServerDeps) *Server {
	return &Server{
		middlewares: deps.Middlewares,

		profile: deps.Profile,

		log:      deps.Log,
		resolver: deps.Resolver,
	}
}

func (s *Server) Run(ctx context.Context, config Config) {
	auth := s.middlewares.AccountAuth()

	r := chi.NewRouter()
	r.Use(
		s.middlewares.Logger(s.log),
		s.middlewares.ResolverUrl(s.resolver),
		s.middlewares.CorsDocs(),
	)

	r.Route("/profiles-svc", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/profiles", func(r chi.Router) {
				r.Get("/", s.profile.Filter)

				r.With(auth).Route("/me", func(r chi.Router) {
					r.Get("/", s.profile.GetMy)
					r.Put("/", s.profile.UpdateMy)

					r.Route("/media", func(r chi.Router) {
						r.Post("/", s.profile.CreateUploadMediaLink)
						r.Delete("/", s.profile.DeleteUploadMedia)
					})
				})

				r.Get("/@{username}", s.profile.GetByUsername)
				r.Get("/{account_id:[0-9a-fA-F-]{36}}", s.profile.GetByID)
			})
		})
	})

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", config.Port),
		Handler:           r,
		ReadTimeout:       config.ReadTimeout,
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		WriteTimeout:      config.WriteTimeout,
		IdleTimeout:       config.IdleTimeout,
	}

	s.log.WithField("port", config.Port).Info("starting http service...")

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
		s.log.Info("shutting down http service...")
	case err := <-errCh:
		if err != nil {
			s.log.WithError(err).Error("http server error")
		}
	}

	shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shCtx); err != nil {
		s.log.WithError(err).Error("failed to shutdown http server gracefully")
	} else {
		s.log.Info("http server shutdown gracefully")
	}
}
