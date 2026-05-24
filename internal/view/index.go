package view

import (
	"net/http"
)

type TemplateManager interface {
	RenderTemplate(w http.ResponseWriter, template_name string, data any)
}

func TemplateView(w http.ResponseWriter, tm TemplateManager, template_name string, data any, status int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	tm.RenderTemplate(w, template_name, data)
}
