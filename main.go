package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/deeerain/nebula-stub/internal/config"
	"github.com/deeerain/nebula-stub/internal/logger"
	"github.com/deeerain/nebula-stub/internal/middlewares"
	"github.com/deeerain/nebula-stub/internal/server"
)

//go:embed frontend/dist/*
var embeddedFiles embed.FS

func Assets() fs.FS {
	distFS, err := fs.Sub(embeddedFiles, "frontend/dist")
	if err != nil {
		panic(err)
	}

	return distFS
}

func main() {
	// Vars
	config := config.FromEnv()
	config.Validate()
	logger.Init("local", slog.LevelDebug)
	logger := slog.Default()
	bindAddress := fmt.Sprintf("%s:%v", config.Listen, config.Port)
	var err error

	frontednFS := Assets()

	server := server.New(logger)

	server.Use(middlewares.LoggingMiddleware(logger))

	// Services
	// templateManager := templates.NewTemplateManager()
	// templateManager.Init(embeddedFS)

	// telemtServic
	// Controllers
	// indexController := controller.NewIndexController(telemtService, templateManager, logger)

	// server.Handle("/static/", http.StripPrefix("/static/", http.FileServerFS(staticFS)))
	server.Handle("/", http.FileServer(http.FS(frontednFS)))

	// Run
	if err = server.Run(bindAddress); err != nil {
		logger.Error("failed to start listener", "error", err)
	}
}
