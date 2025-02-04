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

type HTTPServer struct {
	logger *logger.Logger
	app    *app.App
	server *http.Server
	router *mux.Router
	conf   config.ServerConfig
}

func NewServer(logger *logger.Logger, app *app.App, conf config.ServerConfig) app.Server {
	s := &HTTPServer{
		logger: logger,
		app:    app,
		router: mux.NewRouter(),
		conf:   conf,
	}

	s.setupRoutes()
	return s
}

func (s *HTTPServer) setupRoutes() {
	s.router.Use(func(next http.Handler) http.Handler {
		return LoggingMiddleware(s.logger, next)
	})

	eventHandlers := handlers.NewEventHandlers(s.app, s.logger)

	s.router.HandleFunc("/events", eventHandlers.HandleCreateEvent).Methods("POST")
	s.router.HandleFunc("/events/{id}", eventHandlers.HandleUpdateEvent).Methods("PUT")
	s.router.HandleFunc("/events/{id}", eventHandlers.HandleDeleteEvent).Methods("DELETE")
	s.router.HandleFunc("/events/day",
		eventHandlers.HandleListEventsForDay).Methods("GET").Queries("date", "{date}", "user_id", "{user_id}")
	s.router.HandleFunc("/events/week",
		eventHandlers.HandleListEventsForWeek).Methods("GET").Queries("date", "{date}", "user_id", "{user_id}")
	s.router.HandleFunc("/events/month",
		eventHandlers.HandleListEventsForMonth).Methods("GET").Queries("date", "{date}", "user_id", "{user_id}")
}

func (s *HTTPServer) Start(ctx context.Context) error {
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

func (s *HTTPServer) Stop(ctx context.Context) error {
	s.logger.Info("Stopping HTTP server")
	return s.server.Shutdown(ctx)
}
