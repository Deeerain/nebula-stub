package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"main/internal/controller"
	"main/internal/logger"
	"main/internal/middlewares"
	"main/internal/server"
	"main/internal/service"
	"main/internal/templates"
	"net/http"
	"os"
)

var bindPort = os.Getenv("STUB_PORT")
var TelemtListen string = os.Getenv("STUB_TELEMT_LISTEN")

//go:embed templates/* static/*
var embeddedFS embed.FS

func main() {
	logger.Init("local", slog.LevelDebug)
	logger := slog.Default()

	if bindPort == "" {
		bindPort = "8080"
	}

	staticFS, err := fs.Sub(embeddedFS, "static")
	if err != nil {
		panic(err)
	}

	templateManager := templates.NewTemplateManager()
	templateManager.Init(embeddedFS)

	listen := fmt.Sprintf("0.0.0.0:%v", bindPort)

	server := server.New(logger)

	server.Use(middlewares.LoggingMiddleware(logger))

	telemtService := service.NewTelemtService("localhost", 9091)

	indexController := controller.NewIndexController(telemtService, templateManager, logger)

	server.Handle("/static/", http.StripPrefix("/static/", http.FileServerFS(staticFS)))
	server.HandleFunc("/", indexController.Index)

	if err = server.Run(listen); err != nil {
		logger.Error("failed to start listener", "error", err)
	}
}
