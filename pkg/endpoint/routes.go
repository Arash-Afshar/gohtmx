package endpoint

import "github.com/labstack/echo/v4"

func routes(e *echo.Echo, h Handler) {
	e.GET("/api/sample", h.apiListSampleHandler)
	e.POST("/api/sample", h.apiNewSampleHandler)
	e.DELETE("/api/sample/:name", h.apiDeleteSampleHandler)

	e.GET("/api/posts", h.listPost)
	e.POST("/api/posts", h.createPost)
	e.DELETE("/api/posts/:id", h.deletePost)
}
