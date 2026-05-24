package controller

import (
	"fmt"
	"log"
	"log/slog"
	"main/internal/logger"
	"main/internal/model"
	"main/internal/view"
	"net/http"
)

type TelemtService interface {
	GetConnections() (map[string]any, error)
}

type IndexController struct {
	logger        logger.Logger
	telemtService TelemtService
	tm            view.TemplateManager
}

func NewIndexController(telemtService TelemtService,
	tm view.TemplateManager,
	logger logger.Logger,
) *IndexController {
	if logger == nil {
		logger = slog.Default()
	}

	return &IndexController{
		logger:        logger,
		telemtService: telemtService,
		tm:            tm,
	}
}

func (c *IndexController) Index(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/" {
		view.TemplateView(w, c.tm, "404", nil, 404)
		return
	}
	telemtData, err := c.telemtService.GetConnections()
	if err != nil {
		log.Println(fmt.Errorf("failed to get connections: %w", err))
	}

	pd := model.PageData[map[string]any]{
		Title: "NEBULA",
		Data:  &telemtData,
	}

	view.TemplateView(w, c.tm, "index", pd, 200)
}
