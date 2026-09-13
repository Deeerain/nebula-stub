package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

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

	xuiUrl, _ := url.Parse(os.Getenv("XUI_URL"))
	xuiToken := os.Getenv("XUI_API_TOKEN")
	xrayService := service.NewXrayService(xuiUrl, xuiToken)

	// Services
	// templateManager := templates.NewTemplateManager()
	// templateManager.Init(embeddedFS)

	// telemtServic
	// Controllers
	// indexController := controller.NewIndexController(telemtService, templateManager, logger)

	// server.Handle("/static/", http.StripPrefix("/static/", http.FileServerFS(staticFS)))
	server.Handle("/", http.FileServer(http.FS(frontednFS)))
	server.HandleFunc("/api/status", getContainerStatuses(ctx, dockerService))
	server.HandleFunc("/api/client", getClinet(ctx, xrayService))

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

func getClinet(ctx context.Context, xrayService service.XrayService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorizationValue := r.Header.Get("Authorization")
		initData := strings.Split(authorizationValue, " ")[1]

		botToken := os.Getenv("BOT_TOKEN")
		data, err := ValidateTelegramInitData(botToken, initData, 24*time.Hour)
		if err != nil {
			http.Error(w, err.Error(), 400)
		}

		slog.Info(data.Get("user"))

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}

		w.Header().Set("Content-Type", "application/json")
	}
}

func ValidateTelegramInitData(botToken string, initDataRaw string, expDuration time.Duration) (url.Values, error) {
	queryParams, err := url.ParseQuery(initDataRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse initData: %w", err)
	}

	hashHex := queryParams.Get("hash")
	if hashHex == "" {
		return nil, errors.New("hash is missing in initData")
	}

	if expDuration > 0 {
		authDateStr := queryParams.Get("auth_date")
		if authDateStr == "" {
			return nil, errors.New("auth date is missing")
		}

		authTimestamp, err := strconv.ParseInt(authDateStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid auth_date format: %w", err)
		}

		authTime := time.Unix(authTimestamp, 0)
		if time.Since(authTime) > expDuration {
			return nil, errors.New("initData has expired")
		}
	}

	var keys []string
	for k := range queryParams {
		if k != "hash" {
			keys = append(keys, k)
		}
	}

	var dataPairs []string
	for _, k := range keys {
		dataPairs = append(dataPairs, fmt.Sprintf("%s=%s", k, queryParams.Get(k)))
	}
	dataCheckString := strings.Join(dataPairs, "\n")

	macSeed := hmac.New(sha256.New, []byte("WebAppData"))
	macSeed.Write([]byte(botToken))
	secretKey := macSeed.Sum(nil)

	macHash := hmac.New(sha256.New, secretKey)
	macHash.Write([]byte(dataCheckString))
	calculateHashBytes := macHash.Sum(nil)
	calculateHashHex := hex.EncodeToString(calculateHashBytes)

	if !hmac.Equal([]byte(calculateHashHex), []byte(hashHex)) {
		return nil, errors.New("invalid signature (hash mismatch)")
	}

	return queryParams, nil
}
