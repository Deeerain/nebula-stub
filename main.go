package main

import (
	"context"
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/deeerain/nebula-stub/internal/handlers"
	"github.com/deeerain/nebula-stub/internal/logger"
	"github.com/deeerain/nebula-stub/internal/middlewares"
	"github.com/deeerain/nebula-stub/internal/server"
	"github.com/deeerain/nebula-stub/internal/service"
	"github.com/deeerain/nebula-stub/internal/xui"
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

	slog.SetLogLoggerLevel(slog.LevelDebug)
	// Vars
	logger.Init("local", slog.LevelDebug)
	logger := slog.Default()
	var err error

	ctx := context.Background()

	frontednFS := Assets()

	server := server.New(logger)
	// server.Use(middlewares.CommonMiddleware())
	server.Use(middlewares.LoggingMiddleware(logger))

	dockerService, err := service.NewDockerService()
	if err != nil {
		slog.Warn("docker service is not inited", "error", err)
	}
	defer dockerService.Close()

	xuiUrl, _ := url.Parse(os.Getenv("XUI_URL"))
	xuiToken := os.Getenv("XUI_API_TOKEN")
	xuiClient := xui.NewApiClient(xuiUrl, xuiToken)
	xrayService := xui.NewXuiService(xuiClient)

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	// telegramWebhookUrl := os.Getenv("TELEGRAM_WEBHOOK")

	telegramService, err := service.NewTelegramService(botToken)
	if err != nil {
		slog.Warn("telegram service not inited", "error", err)
	}

	// getSubscriptionUseCase := usecase.NewGetSubscriptionUseCase(xrayService)

	handlers := map[string]http.HandlerFunc{
		"/tgwebhook":              handlers.WebHook(telegramService),
		"/api/status":             handlers.GetContainerStatuses(ctx, dockerService),
		"/api/client":             handlers.GetClinet(ctx, xrayService),
		"/api/subscriptions/{id}": handlers.GetLinks(ctx, xrayService),
		// "/api/client/:id/links": handlers.GetSubscriptionHandler(getSubscriptionUseCase),
	}

	for p, h := range handlers {
		server.HandleFunc(p, h)
	}
	server.HandleFS("/", frontednFS)

	// if err := telegramService.SetWebhook(telegramWebhookUrl); err != nil {
	// 	slog.Warn("telegram webhook not set", "error", err)
	// }

	// Run
	if err = server.Run(":8080"); err != nil {
		logger.Error("failed to start listener", "error", err)
	}

	<-ctx.Done()
}
