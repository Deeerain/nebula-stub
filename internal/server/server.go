package server

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"net"
	"net/http"

	"github.com/deeerain/nebula-stub/internal/logger"
	"github.com/deeerain/nebula-stub/internal/middlewares"
)

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

func (s *Server) Run(listen string) error {
	listenr, err := net.Listen("tcp", listen)
	if err != nil {
		return fmt.Errorf("Failed to create listener: %w", err)
	}

	s.Logger().Info("Server starting", "bind", listenr.Addr())

	chain := middlewares.Chain(s.mux, s.middlewares...)
	return http.Serve(listenr, chain)
}

func WriteJSON(w http.ResponseWriter, obj any) error {
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(obj)
}

func ReadJSON(r *http.Request, obj any) error {
	return json.NewDecoder(r.Body).Decode(obj)
}
