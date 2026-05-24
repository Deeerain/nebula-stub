package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"main/internal/config"
	"main/internal/controller"
	"main/internal/logger"
	"main/internal/middlewares"
	"main/internal/server"
	"main/internal/service"
	"main/internal/templates"
	"net/http"
)

//go:embed templates/* static/*
var embeddedFS embed.FS

func main() {
	// Vars
	config := config.FromEnv()
	config.Validate()
	logger.Init("local", slog.LevelDebug)
	logger := slog.Default()
	bindAddress := fmt.Sprintf("%s:%v", config.Listen, config.Port)
	var staticFS fs.FS
	var err error

	if staticFS, err = fs.Sub(embeddedFS, "static"); err != nil {
		panic(err)
	}

	server := server.New(logger)

	server.Use(middlewares.LoggingMiddleware(logger))

	// Services
	templateManager := templates.NewTemplateManager()
	templateManager.Init(embeddedFS)

	telemtService := service.NewTelemtService(config.TelemtHost, int(config.TelemtPort))

	// Controllers
	indexController := controller.NewIndexController(telemtService, templateManager, logger)

	server.Handle("/static/", http.StripPrefix("/static/", http.FileServerFS(staticFS)))
	server.HandleFunc("/", indexController.Index)

	// Run
	if err = server.Run(bindAddress); err != nil {
		logger.Error("failed to start listener", "error", err)
	}
}
