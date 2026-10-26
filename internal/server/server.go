package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/deeerain/nebula-stub/internal/logger"
	"github.com/deeerain/nebula-stub/internal/middlewares"
)

const shutdownTimeout = 10 * time.Second

type Server struct {
	logger      logger.Logger
	middlewares []middlewares.Middleware
	mux         *http.ServeMux
}

func New(logger logger.Logger) *Server {
	return &Server{
		logger: logger,
		mux:    http.NewServeMux(),
	}
}

func (s *Server) Logger() logger.Logger {
	if s.logger == nil {
		return slog.Default()
	}

	return s.logger
}

func (s *Server) Use(middlewares ...middlewares.Middleware) {
	s.middlewares = append(s.middlewares, middlewares...)
}

func (s *Server) Handle(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

func (s *Server) HandleFunc(pattern string, handler http.HandlerFunc) {
	s.mux.HandleFunc(pattern, handler)
}

func (s *Server) HandleFS(pattern string, fs fs.FS) {
	s.mux.Handle(pattern, http.FileServer(http.FS(fs)))
}

func (s *Server) Run(ctx context.Context, listen string) error {
	listener, err := net.Listen("tcp", listen)
	if err != nil {
		return fmt.Errorf("create listener: %w", err)
	}

	httpServer := &http.Server{
		Handler:           middlewares.Chain(s.mux, s.middlewares...),
		ReadHeaderTimeout: 5 * time.Second,
	}

	serveErr := make(chan error, 1)

	go func() {
		s.Logger().Info("server starting", "bind", listener.Addr().String())

		if err := httpServer.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErr <- err
		}

		close(serveErr)
	}()

	select {
	case err, ok := <-serveErr:
		if ok && err != nil {
			return fmt.Errorf("serve: %w", err)
		}

		return nil
	case <-ctx.Done():
		s.Logger().Info("server shutting down")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown server: %w", err)
		}

		return nil
	}
}

func WriteJSON(w http.ResponseWriter, obj any) error {
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(obj)
}

type errorResponse struct {
	Error string `json:"error"`
}

// WriteError writes a JSON error body and sets the response status.
func WriteError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(errorResponse{Error: message}); err != nil {
		slog.Error("failed to write error response", "error", err)
	}
}

func ReadJSON(r *http.Request, obj any) error {
	return json.NewDecoder(r.Body).Decode(obj)
}
