package endpoint

import (
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
		c.Logger().Errorf("db.ListSamples: err=[%v], method=[%s], status=[%d], path=[%s]", err, c.Request().Method, http.StatusInternalServerError, c.Request().URL.Path)
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
		c.Logger().Errorf("db.AddSample: err=[%v], method=[%s], status=[%d], path=[%s]", err, c.Request().Method, http.StatusInternalServerError, c.Request().URL.Path)
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
		c.Logger().Errorf("db.FindSample: err=[%v], method=[%s], status=[%d], path=[%s]", err, c.Request().Method, http.StatusInternalServerError, c.Request().URL.Path)
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	if err := db.DeleteSample(c.Request().Context(), h.DB, sample); err != nil {
		c.Logger().Errorf("db.DeleteSample: err=[%v], method=[%s], status=[%d], path=[%s]", err, c.Request().Method, http.StatusInternalServerError, c.Request().URL.Path)
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	return h.apiListSampleHandler(c)
}
