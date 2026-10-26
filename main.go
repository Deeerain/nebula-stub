package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/url"
	"os"
	"os/signal"
	"syscall"

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
	if err := run(); err != nil {
		slog.Error("application stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	logger.Init(slog.LevelDebug)
	log := slog.Default()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	frontendFS := Assets()

	xuiURL, err := absoluteURL(env("XUI_URL"))
	if err != nil {
		return fmt.Errorf("invalid XUI_URL: %w", err)
	}

	subBaseURL, err := absoluteURL(env("SUBSCRIPTION_BASE_URL"))
	if err != nil {
		return fmt.Errorf("invalid SUBSCRIPTION_BASE_URL: %w", err)
	}

	dockerService, err := service.NewDockerService()
	if err != nil {
		log.Warn("docker service is not initialized", "error", err)
	}
	if dockerService != nil {
		defer dockerService.Close()
	}

	xuiClient := xui.NewApiClient(xuiURL, env("XUI_API_TOKEN", "XRAY_API_TOKEN"))
	xrayService := xui.NewXuiService(xuiClient)

	botToken := env("TELEGRAM_BOT_TOKEN", "BOT_TOKEN")

	telegramService, err := service.NewTelegramService(botToken)
	if err != nil {
		log.Warn("telegram service is not initialized", "error", err)
	}

	useAuth := auth.AuthMiddleware()

	srv := server.New(log)
	srv.Use(middlewares.LoggingMiddleware(log))

	srv.Handle("/api/client", useAuth(client.GetClient(ctx, xrayService)))
	srv.Handle("/api/connect/{app}", useAuth(subscription.ConnectHandler(*subBaseURL, xrayService)))
	srv.Handle("/api/subscriptions/{id}", useAuth(subscription.GetLinks(ctx, xrayService)))
	// srv.Handle("/api/payments/invoice", useAuth(payments.InvoiceHandler(paymentService)))

	devMode := os.Getenv("STUB_DEV") == "1"
	if devMode {
		log.Warn(`STUB_DEV is enabled: /api/auth accepts the "dev" initData`)
	}

	srv.HandleFunc("/api/auth", auth.AuthHandler(botToken, devMode))
	srv.HandleFunc("/tgwebhook", handlers.WebHook(telegramService))
	srv.HandleFunc("/api/status", handlers.GetContainerStatuses(dockerService))
	srv.HandleFS("/", frontendFS)

	if err := srv.Run(ctx, ":8080"); err != nil {
		return fmt.Errorf("run server: %w", err)
	}

	return nil
}

func absoluteURL(raw string) (*url.URL, error) {
	if raw == "" {
		return nil, errors.New("value is not set")
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("parse %q: %w", raw, err)
	}

	if parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("must be an absolute URL, got %q", raw)
	}

	return parsed, nil
}

func env(keys ...string) string {
	for _, key := range keys {
		if value := os.Getenv(key); value != "" {
			return value
		}
	}

	return ""
}
