package endpoint

import (
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"

	dbPkg "github.com/Arash-Afshar/gohtmx/pkg/db"
	"github.com/Arash-Afshar/gohtmx/pkg/models"
	"github.com/go-chi/chi/v5"
)

func ParsePartialTemplate(names ...string) *template.Template {
	pagePath := filepath.Join(names...)
	return template.Must(template.ParseFiles(pagePath))
}

func ParseTemplate(names ...string) *template.Template {
	global := []string{
		filepath.Join("templates", "main.html"),
		filepath.Join("templates", "partials", "samples.html"), // TODO: glob
	}

	// Add all user templates after global.
	pagePath := filepath.Join(names...)
	global = append(global, pagePath)

	return template.Must(template.ParseFiles(global...))
}

func ExecuteTemplate(w http.ResponseWriter, r *http.Request, data any, names ...string) bool {
	tmpl := ParseTemplate(names...)

	if err := tmpl.Execute(w, data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error(err.Error(), "method", r.Method, "status", http.StatusInternalServerError, "path", r.URL.Path)
		return false
	}
	return true
}

func ExecutePartialTemplate(w http.ResponseWriter, r *http.Request, data any, names ...string) bool {
	tmpl := ParsePartialTemplate(names...)

	if err := tmpl.Execute(w, data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error(err.Error(), "method", r.Method, "status", http.StatusInternalServerError, "path", r.URL.Path)
		return false
	}
	return true
}

func isHtmx(w http.ResponseWriter, r *http.Request) bool {
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

func (h *Handler) indexViewHandler(w http.ResponseWriter, r *http.Request) {
	if !ExecuteTemplate(w, r, nil, "templates", "pages", "index.html") {
		return
	}
	slog.Info("render page", "method", r.Method, "status", http.StatusOK, "path", r.URL.Path)
}

func (h *Handler) apiListSampleHandler(w http.ResponseWriter, r *http.Request) {
	if !isHtmx(w, r) {
		return
	}
	samples, err := dbPkg.ListSamples(h.DB)
	if err != nil {
		slog.Error("Getting samples list", "errMessage", err, "method", r.Method, "status", http.StatusInternalServerError, "path", r.URL.Path)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	type data struct {
		Samples []*models.Sample
	}
	if !ExecutePartialTemplate(w, r, data{Samples: samples}, "templates", "partials", "samples.html") {
		return
	}
	slog.Info("render page", "method", r.Method, "status", http.StatusOK, "path", r.URL.Path)
}

func (h *Handler) apiNewSampleHandler(w http.ResponseWriter, r *http.Request) {
	if !isHtmx(w, r) {
		return
	}
	name := r.FormValue("name")
	sample := models.NewSample(name)
	if err := dbPkg.AddSample(h.DB, sample); err != nil {
		slog.Error("Add a sample", "errMessage", err, "method", r.Method, "status", http.StatusInternalServerError, "path", r.URL.Path)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	h.apiListSampleHandler(w, r)
}

func (h *Handler) apiDeleteSampleHandler(w http.ResponseWriter, r *http.Request) {
	if !isHtmx(w, r) {
		return
	}
	name := chi.URLParam(r, "name")
	sample, err := dbPkg.SampleByName(h.DB, name)
	if err != nil {
		slog.Error("Find the sample", "errMessage", err, "method", r.Method, "status", http.StatusInternalServerError, "path", r.URL.Path)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if err := dbPkg.DeleteSample(h.DB, sample); err != nil {
		slog.Error("Delete a sample", "errMessage", err, "method", r.Method, "status", http.StatusInternalServerError, "path", r.URL.Path)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	h.apiListSampleHandler(w, r)
}
