package endpoint

import (
	"log/slog"
	"net/http"

	"github.com/Arash-Afshar/gohtmx/pkg/db"
	"github.com/Arash-Afshar/gohtmx/pkg/models"
	"github.com/labstack/echo/v4"
)

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

func (h *Handler) indexViewHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "pages/index.html", nil)
}

func (h *Handler) apiListSampleHandler(c echo.Context) error {
	samples, err := db.ListSamples(h.DB)
	if err != nil {
		slog.Error("Getting samples list", "errMessage", err, "method", c.Request().Method, "status", http.StatusInternalServerError, "path", c.Request().URL.Path)
		return c.Render(http.StatusInternalServerError, "pages/error.html", nil)
	}
	type data struct {
		Samples []*models.Sample
	}
	return c.Render(http.StatusOK, "partials/samples.html", data{Samples: samples})
}

func (h *Handler) apiNewSampleHandler(c echo.Context) error {
	if !isHtmx(c) {
		return nil
	}
	name := c.FormValue("name")
	sample := models.NewSample(name)
	if err := db.AddSample(h.DB, sample); err != nil {
		slog.Error("Add a sample", "errMessage", err, "method", c.Request().Method, "status", http.StatusInternalServerError, "path", c.Request().URL.Path)
		return c.Render(http.StatusInternalServerError, "pages/error.html", nil)
	}
	return h.apiListSampleHandler(c)
}

func (h *Handler) apiDeleteSampleHandler(c echo.Context) error {
	if !isHtmx(c) {
		return nil
	}
	name := c.Param("name")
	sample, err := db.SampleByName(h.DB, name)
	if err != nil {
		slog.Error("Find the sample", "errMessage", err, "method", c.Request().Method, "status", http.StatusInternalServerError, "path", c.Request().URL.Path)
		return c.Render(http.StatusInternalServerError, "pages/error.html", nil)
	}
	if err := db.DeleteSample(h.DB, sample); err != nil {
		slog.Error("Delete a sample", "errMessage", err, "method", c.Request().Method, "status", http.StatusInternalServerError, "path", c.Request().URL.Path)
		return c.Render(http.StatusInternalServerError, "pages/error.html", nil)
	}
	return h.apiListSampleHandler(c)
}
