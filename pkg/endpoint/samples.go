package endpoint

import (
	"log/slog"
	"net/http"

	"github.com/Arash-Afshar/gohtmx/pkg/db"
	"github.com/Arash-Afshar/gohtmx/pkg/models"
	"github.com/labstack/echo/v4"
)

func (h *Handler) indexViewHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "pages/samples.html", nil)
}

func (h *Handler) apiListSampleHandler(c echo.Context) error {
	samples, err := db.ListSamples(c.Request().Context(), h.DB)
	if err != nil {
		slog.Error("Getting samples list", "errMessage", err, "method", c.Request().Method, "status", http.StatusInternalServerError, "path", c.Request().URL.Path)
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	type data struct {
		Samples []models.Sample
	}
	return c.Render(http.StatusOK, "partials/samples.html", data{Samples: samples})
}

func (h *Handler) apiNewSampleHandler(c echo.Context) error {
	if !isHtmx(c) {
		return nil
	}
	name := c.FormValue("name")
	sample := models.NewSample(name)
	if err := db.AddSample(c.Request().Context(), h.DB, sample); err != nil {
		slog.Error("Add a sample", "errMessage", err, "method", c.Request().Method, "status", http.StatusInternalServerError, "path", c.Request().URL.Path)
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	return h.apiListSampleHandler(c)
}

func (h *Handler) apiDeleteSampleHandler(c echo.Context) error {
	if !isHtmx(c) {
		return nil
	}
	id := c.Param("id")
	sample, err := db.FindSample(c.Request().Context(), h.DB, id)
	if err != nil {
		slog.Error("Find the sample", "errMessage", err, "method", c.Request().Method, "status", http.StatusInternalServerError, "path", c.Request().URL.Path)
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	if err := db.DeleteSample(c.Request().Context(), h.DB, sample); err != nil {
		slog.Error("Delete a sample", "errMessage", err, "method", c.Request().Method, "status", http.StatusInternalServerError, "path", c.Request().URL.Path)
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	return h.apiListSampleHandler(c)
}
