package endpoint

import (
	"log/slog"
	"net/http"

	"github.com/Arash-Afshar/gohtmx/pkg/db"
	"github.com/Arash-Afshar/gohtmx/pkg/models"
	"github.com/labstack/echo/v4"
)

func (h *Handler) listPost(c echo.Context) error {
	posts, err := db.ListPosts(c.Request().Context(), h.DB)
	if err != nil {
		slog.Error("Getting samples list", "errMessage", err, "method", c.Request().Method, "status", http.StatusInternalServerError, "path", c.Request().URL.Path)
		return c.Render(http.StatusInternalServerError, "pages/error.html", DisplayError{Message: err.Error()})
	}
	type data struct {
		Posts []models.Post
	}
	return c.Render(http.StatusOK, "pages/posts.html", data{Posts: posts})
}

func (h *Handler) createPost(c echo.Context) error {
	if !isHtmx(c) {
		return nil
	}
	title := c.FormValue("title")
	post := models.NewPost(title)
	if err := db.AddPost(c.Request().Context(), h.DB, post); err != nil {
		slog.Error("createPost", "err", err, "method", c.Request().Method, "status", http.StatusInternalServerError, "path", c.Request().URL.Path)
		return c.Render(http.StatusInternalServerError, "pages/error.html", DisplayError{Message: err.Error()})
	}
	return h.listPost(c)
}

func (h *Handler) deletePost(c echo.Context) error {
	if !isHtmx(c) {
		return nil
	}
	id := c.Param("id")
	post, err := db.FindPost(c.Request().Context(), h.DB, id)
	if err != nil {
		slog.Error("Find the sample", "errMessage", err, "method", c.Request().Method, "status", http.StatusInternalServerError, "path", c.Request().URL.Path)
		return c.Render(http.StatusInternalServerError, "pages/error.html", nil)
	}
	if err := db.DeletePost(c.Request().Context(), h.DB, post); err != nil {
		slog.Error("Delete a sample", "errMessage", err, "method", c.Request().Method, "status", http.StatusInternalServerError, "path", c.Request().URL.Path)
		return c.Render(http.StatusInternalServerError, "pages/error.html", nil)
	}
	return h.listPost(c)
}
