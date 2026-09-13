package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net/url"
	"os"

	"github.com/deeerain/nebula-stub/internal/config"
	"github.com/deeerain/nebula-stub/internal/handlers"
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
	server.Use(middlewares.CommonMiddleware())
	server.Use(middlewares.LoggingMiddleware(logger))

	dockerService, err := service.NewDockerService()
	if err != nil {
		slog.Warn("docker service is not inited", "error", err)
	}
	defer dockerService.Close()

	xuiUrl, _ := url.Parse(os.Getenv("XUI_URL"))
	xuiToken := os.Getenv("XUI_API_TOKEN")
	xrayService := service.NewXrayService(xuiUrl, xuiToken)

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	telegramService, err := service.NewTelegramService(botToken)
	if err != nil {
		slog.Warn("telegram service not inited", "error", err)
	}

	server.HandleFS("/", frontednFS)
	server.HandleFunc("/tgwebhook", handlers.WebHook(telegramService))
	server.HandleFunc("/api/status", handlers.GetContainerStatuses(ctx, dockerService))
	server.HandleFunc("/api/client", handlers.GetClinet(ctx, xrayService))

	// Run
	if err = server.Run(bindAddress); err != nil {
		logger.Error("failed to start listener", "error", err)
	}

	<-ctx.Done()
}
