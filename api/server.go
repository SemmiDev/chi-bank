package api

import (
	"fmt"
	"net/http"

	db "github.com/SemmiDev/chi-bank/db/sqlc"
	"github.com/SemmiDev/chi-bank/util/config"
	"github.com/SemmiDev/chi-bank/util/logger"
	"github.com/SemmiDev/chi-bank/util/token"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config     config.Env
	store      db.Store
	tokenMaker token.Maker
	logger     *logger.Logger
	router     *chi.Mux
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config config.Env, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	lg := logger.New(true)

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		logger:     lg,
	}

	server.setupRouter()
	return server, nil
}

func (s *Server) setupRouter() {
	r := chi.NewRouter()

	r.Use(cors.AllowAll().Handler)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Options("/*", func(w http.ResponseWriter, r *http.Request) {})
	// api := r.Route("/api/v1", func(router chi.Router) {})

	// api.Route("/accounts", func(r chi.Router) {
	// 	r.Post("/users", s.createUserHandler()),
	// }

	s.router = r
}

// Start runs the HTTP server on a specific address.
func (s *Server) Start(address string) error {
	err := http.ListenAndServe(address, s.router)
	if err != nil {
		return err
	}
	return nil
}
