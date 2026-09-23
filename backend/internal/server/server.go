package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"lumi.yellowlabs.space/internal/api/http/routes"
	"lumi.yellowlabs.space/internal/config"
)

type Server struct {
	router chi.Router
	server *http.Server
}

func NewServer() *Server {
	router := chi.NewRouter()
	routes.SetupRoutes(router)

	return &Server{
		router: router,
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", config.AppConfig.Server.Host, config.AppConfig.Server.Port),
			Handler: router,
		},
	}
}

func (s *Server) Start() error {
	log.Infof("Starting server on %s", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}
