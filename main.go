package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/deeerain/nebula-stub/internal/config"
	"github.com/deeerain/nebula-stub/internal/logger"
	"github.com/deeerain/nebula-stub/internal/middlewares"
	"github.com/deeerain/nebula-stub/internal/server"
	"github.com/deeerain/nebula-stub/internal/service"
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

	ctx := context.Background()

	frontednFS := Assets()

	server := server.New(logger)

	server.Use(middlewares.LoggingMiddleware(logger))

	dockerService, err := service.NewDockerService()
	if err != nil {
		slog.Warn("docker service is not inited", "error", err)
	}
	defer dockerService.Close()

	// Services
	// templateManager := templates.NewTemplateManager()
	// templateManager.Init(embeddedFS)

	// telemtServic
	// Controllers
	// indexController := controller.NewIndexController(telemtService, templateManager, logger)

	// server.Handle("/static/", http.StripPrefix("/static/", http.FileServerFS(staticFS)))
	server.Handle("/", http.FileServer(http.FS(frontednFS)))
	server.HandleFunc("/api/status", getContainerStatuses(ctx, dockerService))

	// Run
	if err = server.Run(bindAddress); err != nil {
		logger.Error("failed to start listener", "error", err)
	}

	<-ctx.Done()
}

func getContainerStatuses(ctx context.Context, dockerService service.DockerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dockerService == nil {
			http.NotFound(w, r)
		}

		containers, err := dockerService.GetContainerList(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(containers)
	}
}
