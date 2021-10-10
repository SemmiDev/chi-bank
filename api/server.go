package api

import (
	"context"
	"fmt"
	"github.com/SemmiDev/chi-bank/common"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/SemmiDev/chi-bank/common/logger"
	"github.com/SemmiDev/chi-bank/common/token"
	db "github.com/SemmiDev/chi-bank/db/sqlc"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	//"github.com/go-chi/httprate"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config     common.Config
	store      db.Store
	tokenMaker token.Maker
	logger     *logger.Logger
	router     *chi.Mux
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config common.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	log := logger.New(true)
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		logger:     log,
	}
	server.setupRouter()
	return server, nil
}

func (s *Server) setupRouter() {
	r := chi.NewRouter()

	//r.Use(httprate.LimitByIP(
	//	config.Cfg().HttpRateLimitRequest,
	//	config.Cfg().HttpRateLimitTime,
	//))
	r.Use(cors.AllowAll().Handler)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		// users route
		r.Route("/users", func(r chi.Router) {
			r.Post("/", s.CreateUserHandler)
			r.Post("/login", s.LoginUserHandler)
		})

		// accounts route
		r.With(s.AuthMiddleware).Route("/accounts", func(r chi.Router) {
			r.Post("/", s.createAccount)
			r.Get("/", s.listAccounts)
			r.Get("/{id}", s.getAccount)
		})

		// transfer route
		r.With(s.AuthMiddleware).Post("/transfers", s.createTransfer)
	})
	s.router = r
}

// Start runs the HTTP server on a specific address.
func (s *Server) Start(address string) error {
	httpServer := &http.Server{
		Addr:    address,
		Handler: s.router,
	}

	idleConsClosed := make(chan struct{})
	go func() {
		defer close(idleConsClosed)

		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)

		<-sigint

		err := httpServer.Shutdown(context.Background())
		if err != nil {
			s.logger.Error().Msg("failed to shutdown server")
		}
	}()

	s.logger.Info().Msgf("starting server on %s", httpServer.Addr)
	err := httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	<-idleConsClosed

	s.logger.Info().Msg("stopped server gracefully")
	return nil
}
