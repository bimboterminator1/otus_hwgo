package internalhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"

	//nolint:depguard
	app "github.com/bimboterminator1/otus_hwgo/hw12_13_14_15_calendar/internal/app"
	"github.com/bimboterminator1/otus_hwgo/hw12_13_14_15_calendar/internal/config"
	"github.com/bimboterminator1/otus_hwgo/hw12_13_14_15_calendar/internal/logger"
	"github.com/bimboterminator1/otus_hwgo/hw12_13_14_15_calendar/internal/server/http/handlers"
	"github.com/gorilla/mux"
)

type Server struct {
	logger *logger.Logger
	app    app.App
	server *http.Server
	router *mux.Router
	conf   config.ServerConfig
}

func NewServer(logger *logger.Logger, app app.App, conf config.ServerConfig) *Server {
	s := &Server{
		logger: logger,
		app:    app,
		router: mux.NewRouter(),
		conf:   conf,
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.Use(func(next http.Handler) http.Handler {
		return LoggingMiddleware(s.logger, next)
	})
	s.router.HandleFunc("/hello", handlers.HandleHello).Methods("GET")

	// s.router.HandleFunc("/events/{id}", s.handleDeleteEvent).Methods("DELETE")

	// s.router.HandleFunc("/events", s.handleListEvents).Methods("GET")
}

func (s *Server) Start(ctx context.Context) error {
	addr := net.JoinHostPort(s.conf.Listener.Host, strconv.Itoa(s.conf.Listener.Port))
	s.server = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	go func() {
		<-ctx.Done()
		s.Stop(context.Background())
	}()

	s.logger.Info(fmt.Sprintf("Starting HTTP server on %s", addr))
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping HTTP server")
	return s.server.Shutdown(ctx)
}
