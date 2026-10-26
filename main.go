package main

import (
	"context"
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/deeerain/nebula-stub/internal/auth"
	"github.com/deeerain/nebula-stub/internal/client"
	"github.com/deeerain/nebula-stub/internal/handlers"
	"github.com/deeerain/nebula-stub/internal/logger"
	"github.com/deeerain/nebula-stub/internal/middlewares"
	"github.com/deeerain/nebula-stub/internal/server"
	"github.com/deeerain/nebula-stub/internal/service"
	"github.com/deeerain/nebula-stub/internal/subscription"
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
	logger.Init("local", slog.LevelDebug)
	logger := slog.Default()
	var err error

	ctx := context.Background()

	frontednFS := Assets()

	useAuth := auth.AuthMiddleware()

	server := server.New(logger)
	server.Use(middlewares.LoggingMiddleware(logger))

	dockerService, err := service.NewDockerService()
	if err != nil {
		slog.Warn("docker service is not initialized", "error", err)
	}
	if dockerService != nil {
		defer dockerService.Close()
	}

	xuiURL, err := url.Parse(env("XUI_URL"))
	if err != nil || xuiURL.Host == "" {
		slog.Error(
			"XUI_URL must be an absolute URL, for example https://panel.example.com",
			"value", os.Getenv("XUI_URL"),
		)
	}

	xuiToken := env("XUI_API_TOKEN", "XRAY_API_TOKEN")
	xuiClient := xui.NewApiClient(xuiURL, xuiToken)
	xrayService := xui.NewXuiService(xuiClient)

	botToken := env("TELEGRAM_BOT_TOKEN", "BOT_TOKEN")
	// telegramWebhookUrl := os.Getenv("TELEGRAM_WEBHOOK")

	telegramService, err := service.NewTelegramService(botToken)
	if err != nil {
		slog.Warn("telegram service not inited", "error", err)
	}

	subEndpoint := url.URL{
		Scheme: "https",
		Host:   "nebulasrv.ignorelist.com:8822",
		Path:   "/target-subscribes"}

	server.Handle("/api/client", useAuth(client.GetClient(ctx, xrayService)))
	server.Handle("/api/connect/{app}", useAuth(subscription.ConnectHandler(subEndpoint, xrayService)))
	server.Handle("/api/subscriptions/{id}", useAuth(subscription.GetLinks(ctx, xrayService)))

	devMode := os.Getenv("STUB_DEV") == "1"
	if devMode {
		slog.Warn(`STUB_DEV is enabled: /api/auth accepts the "dev" initData`)
	}

	handlers := map[string]http.HandlerFunc{
		"/api/auth":   auth.AuthHandler(botToken, devMode),
		"/tgwebhook":  handlers.WebHook(telegramService),
		"/api/status": handlers.GetContainerStatuses(dockerService),
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

func env(keys ...string) string {
	for _, key := range keys {
		if value := os.Getenv(key); value != "" {
			return value
		}
	}

	return ""
}
