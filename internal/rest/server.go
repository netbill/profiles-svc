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

	UpdateMyProfile(w http.ResponseWriter, r *http.Request)
	UpdateProfileOfficial(w http.ResponseWriter, r *http.Request)

	CreateMyProfileUploadMediaLink(w http.ResponseWriter, r *http.Request)
	DeleteMyProfileUploadAvatar(w http.ResponseWriter, r *http.Request)
}

type Middlewares interface {
	AccountAuth(
		allowedRoles ...string,
	) func(next http.Handler) http.Handler
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
	Port     string `mapstructure:"port" required:"true"`
	Timeouts struct {
		Read       time.Duration `mapstructure:"read"`
		ReadHeader time.Duration `mapstructure:"read_header"`
		Write      time.Duration `mapstructure:"write"`
		Idle       time.Duration `mapstructure:"idle"`
	} `mapstructure:"timeouts"`
}

func (s *Server) Run(ctx context.Context, log *logium.Entry, cfg Config) {
	auth := s.middlewares.AccountAuth()
	sysmoder := s.middlewares.AccountAuth(tokens.RoleSystemAdmin, tokens.RoleSystemModer)

	log = log.WithField("component", "rest")

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
		Addr:              cfg.Port,
		Handler:           r,
		ReadTimeout:       cfg.Timeouts.Read,
		ReadHeaderTimeout: cfg.Timeouts.ReadHeader,
		WriteTimeout:      cfg.Timeouts.Write,
		IdleTimeout:       cfg.Timeouts.Idle,
	}

	log.Infof("starting http service on %s", cfg.Port)

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
			log.Errorf("http server error: %v", err)
		}
	}

	shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shCtx); err != nil {
		log.Errorf("http shutdown error: %v", err)
	} else {
		log.Infof("REST server stopped")
	}
}
