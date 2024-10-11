package endpoint

import (
	"net/http"

	"github.com/Arash-Afshar/gohtmx/pkg/db"
	"github.com/Arash-Afshar/gohtmx/pkg/models"
	"github.com/labstack/echo/v4"
)

type CreatePost struct {
	Title   string `param:"title" query:"title" form:"title" validate:"required"`
	Content string `param:"content" query:"content" form:"content"`
}

type DeletePost struct {
	ID string `param:"id" query:"id" form:"id"`
}

func (h *Handler) listPost(c echo.Context) error {
	posts, err := db.ListPosts(c.Request().Context(), h.DB)
	if err != nil {
		c.Logger().Errorf("db.listPosts: err=[%v], method=[%s], status=[%d], path=[%s]", err, c.Request().Method, http.StatusInternalServerError, c.Request().URL.Path)
		return echo.NewHTTPError(http.StatusUnauthorized, err)
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
	var p CreatePost
	err := c.Bind(&p)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	if err = c.Validate(p); err != nil {
		return err
	}
	post := models.NewPost(p.Title, p.Content)
	if err := db.AddPost(c.Request().Context(), h.DB, post); err != nil {
		c.Logger().Errorf("db.AddPost: err=[%v], method=[%s], status=[%d], path=[%s]", err, c.Request().Method, http.StatusInternalServerError, c.Request().URL.Path)
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	return h.listPost(c)
}

func (h *Handler) deletePost(c echo.Context) error {
	if !isHtmx(c) {
		return nil
	}
	var p DeletePost
	err := c.Bind(&p)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	post, err := db.FindPost(c.Request().Context(), h.DB, p.ID)
	if err != nil {
		c.Logger().Errorf("db.FindPost: err=[%v], method=[%s], status=[%d], path=[%s]", err, c.Request().Method, http.StatusInternalServerError, c.Request().URL.Path)
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	if err := db.DeletePost(c.Request().Context(), h.DB, post); err != nil {
		c.Logger().Errorf("db.DeletePost: err=[%v], method=[%s], status=[%d], path=[%s]", err, c.Request().Method, http.StatusInternalServerError, c.Request().URL.Path)
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	return h.listPost(c)
}
