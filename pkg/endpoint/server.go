package endpoint

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"

	"github.com/Arash-Afshar/gohtmx/pkg/db"
	"github.com/caarlos0/env/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Handler struct {
	DB *sql.DB
}

type Config struct {
	Address string `env:"GOHTMX_BACKEND_HOST" envDefault:"127.0.0.1:7000"`
	DbURL   string `env:"GOHTMX_BACKEND_DB_URL" envDefault:"db.sqlite"`
}

type Templates struct {
	templates map[string]*template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		err := errors.New("Template not found -> " + name)
		return err
	}
	return tmpl.Execute(w, data)
}

func NewTemplates() *Templates {
	mainPath := "templates/main.html"
	pagesGlob := "templates/pages/*.html"
	partialsGlob := "templates/partials/*.html"

	// find matched files
	pages, err := filepath.Glob(pagesGlob)
	if err != nil {
		panic(err)
	}
	partials, err := filepath.Glob(partialsGlob)
	if err != nil {
		panic(err)
	}

	templates := make(map[string]*template.Template, 0)
	for _, page := range pages {
		files := []string{mainPath, page}
		files = append(files, partials...)
		t, err := template.ParseFiles(files...)
		if err != nil {
			panic(err)
		}
		templates["pages/"+filepath.Base(page)] = t
	}

	for _, partial := range partials {
		t, err := template.ParseFiles(partial)
		if err != nil {
			panic(err)
		}
		templates["partials/"+filepath.Base(partial)] = t
	}

	return &Templates{
		templates: templates,
	}
}

type DisplayError struct {
	Message string
}

func Run() error {
	config := Config{}
	if err := env.Parse(&config); err != nil {
		panic(err)
	}

	e := echo.New()
	e.Renderer = NewTemplates()
	e.Use(middleware.Logger())

	var err error
	dbInstance, err := db.NewDB(config.DbURL)
	if err != nil {
		return fmt.Errorf("db connection: %v", err)
	}

	h := Handler{DB: dbInstance}

	// Handle pages
	e.Static("/static", "static")
	routes(e, h)

	e.Logger.Fatal(e.Start(config.Address))
	return nil
}

func isHtmx(c echo.Context) bool {
	w := c.Response()
	r := c.Request()
	// Check, if the current request has a 'HX-Request' header.
	// For more information, see https://htmx.org/docs/#request-headers
	if r.Header.Get("HX-Request") == "" || r.Header.Get("HX-Request") != "true" {
		// If not, return HTTP 400 error.
		w.WriteHeader(http.StatusBadRequest)
		slog.Error("request API", "method", r.Method, "status", http.StatusBadRequest, "path", r.URL.Path)
		return false
	}
	return true
}
