package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
)

var bindPort = os.Getenv("STUB_PORT")
var telemtListen string = os.Getenv("STUB_TELEMT_LISTEN")

//go:embed templates/* static/*
var embeddedFS embed.FS

var templates = make(map[string]*template.Template)

type PageData[T any] struct {
	Title string
	Data  *T
}

func main() {
	if bindPort == "" {
		bindPort = "8080"
	}

	mux := http.NewServeMux()

	staticFS, err := fs.Sub(embeddedFS, "static")
	if err != nil {
		panic(err)
	}

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServerFS(staticFS)))

	mux.HandleFunc("/", homeHandler())

	log.Println("Server wil listen on: ", bindPort)
	log.Println("Telemt api address: ", telemtListen)

	if err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", bindPort), mux); err != nil {
		log.Panicln("failed to start listener: %w", err)
	}
}

func init() {
	templates["index"] = template.Must(
		template.New("base.html").Funcs(template.FuncMap{
			"safeURL": func(s string) template.URL {
				return template.URL(s)
			},
		}).ParseFS(embeddedFS, "templates/base.html", "templates/index.html"),
	)
}

func renderTemplate(w http.ResponseWriter, templName string, data any) {
	temp, exists := templates[templName]
	if !exists {
		http.Error(w, "template not found", http.StatusInternalServerError)
	}

	w.Header().Set("Content type", "text/html; charset=utf-8")

	err := temp.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func homeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		telemtData, err := telemtGetConnections(http.DefaultClient, telemtListen)
		if err != nil {
			log.Println(fmt.Errorf("failed to get connections: %w", err))
		}

		pd := PageData[map[string]any]{
			Title: "NEBULA",
			Data:  &telemtData,
		}
		renderTemplate(w, "index", pd)
	}
}

func telemtGetConnections(client *http.Client, telemtListen string) (map[string]any, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/v1/users", telemtListen), nil)
	if err != nil {
		log.Fatalln("failed to create request: ", err.Error())
	}
	req.Header.Add("Content-type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}

	var data map[string]any
	if err = json.NewDecoder(response.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode body: %w", err)
	}

	return data, nil
}
