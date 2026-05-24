package templates

import (
	"html/template"
	"io/fs"
	"net/http"
)

type TemplateManager struct {
	templates map[string]*template.Template
}

func NewTemplateManager() *TemplateManager {
	return &TemplateManager{
		templates: make(map[string]*template.Template),
	}
}

func (tm *TemplateManager) Init(fs fs.FS) {
	tm.templates["index"] = template.Must(
		template.New("base.html").Funcs(template.FuncMap{
			"safeURL": func(s string) template.URL {
				return template.URL(s)
			},
		}).ParseFS(fs, "templates/*.html"),
	)
}

func (tm *TemplateManager) RenderTemplate(w http.ResponseWriter, templName string, data any) {
	temp, exists := tm.templates[templName]
	if !exists {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content type", "text/html; charset=utf-8")

	err := temp.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
