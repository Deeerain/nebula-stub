package main

import (
	"fmt"
	"log/slog"
	"net/http"

	nebulastub "github.com/deeerain/nebula-stub"
	"github.com/deeerain/nebula-stub/internal/config"
	"github.com/deeerain/nebula-stub/internal/logger"
	"github.com/deeerain/nebula-stub/internal/middlewares"
	"github.com/deeerain/nebula-stub/internal/server"
)

func main() {
	// Vars
	config := config.FromEnv()
	config.Validate()
	logger.Init("local", slog.LevelDebug)
	logger := slog.Default()
	bindAddress := fmt.Sprintf("%s:%v", config.Listen, config.Port)
	var err error

	frontednFS := nebulastub.Asstest()

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
