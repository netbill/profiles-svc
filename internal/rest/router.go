package rest

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/netbill/logium"
	"github.com/netbill/restkit/tokens"
)

type Handlers interface {
	GetMyProfile(w http.ResponseWriter, r *http.Request)

	GetProfileByUsername(w http.ResponseWriter, r *http.Request)
	GetProfileByID(w http.ResponseWriter, r *http.Request)

	FilterProfiles(w http.ResponseWriter, r *http.Request)

	UpdateProfileOfficial(w http.ResponseWriter, r *http.Request)

	OenUpdateProfileSession(w http.ResponseWriter, r *http.Request)
	ConfirmUpdateMyProfile(w http.ResponseWriter, r *http.Request)
	CancelUpdateProfileSession(w http.ResponseWriter, r *http.Request)
	DeleteUploadProfileAvatar(w http.ResponseWriter, r *http.Request)
}

type Middlewares interface {
	AccountAuth(
		allowedRoles ...string,
	) func(next http.Handler) http.Handler
	UpdateOwnProfileMediaContent() func(next http.Handler) http.Handler
	Logger(log *logium.Entry) func(next http.Handler) http.Handler
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
	Port              string
	TimeoutRead       time.Duration
	TimeoutReadHeader time.Duration
	TimeoutWrite      time.Duration
	TimeoutIdle       time.Duration
}

func (s *Server) Run(ctx context.Context, log *logium.Entry, cfg Config) {
	auth := s.middlewares.AccountAuth()
	sysmoder := s.middlewares.AccountAuth(tokens.RoleSystemAdmin, tokens.RoleSystemModer)
	updateOwnProfile := s.middlewares.UpdateOwnProfileMediaContent()

	r := chi.NewRouter()

	r.Use(
		s.middlewares.Logger(log),
		s.middlewares.CorsDocs(),
	)

	r.Route("/profiles-svc", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/profiles", func(r chi.Router) {
				r.Get("/", s.handlers.FilterProfiles)

				r.Get("/u/{username}", s.handlers.GetProfileByUsername)

				r.With(auth).Route("/me", func(r chi.Router) {
					r.Get("/", s.handlers.GetMyProfile)

					r.Route("/update-session", func(r chi.Router) {
						r.Post("/", s.handlers.OenUpdateProfileSession)
						r.With(updateOwnProfile).Delete("/", s.handlers.CancelUpdateProfileSession)

						r.With(updateOwnProfile).Put("/confirm", s.handlers.ConfirmUpdateMyProfile)
						r.With(updateOwnProfile).Delete("/upload-avatar", s.handlers.DeleteUploadProfileAvatar)
					})
				})
			})

			r.Route("/{account_id}", func(r chi.Router) {
				r.Get("/", s.handlers.GetProfileByID)

				r.With(sysmoder).Patch("/official", s.handlers.UpdateProfileOfficial)
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

	log.Infof("starting REST service on %s", cfg.Port)

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
		log.Warnf("shutting down REST service...")
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
		log.Warnf("REST server stopped")
	}
}
