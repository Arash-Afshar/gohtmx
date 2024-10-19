package endpoint

import "github.com/labstack/echo/v4"

func routes(e *echo.Echo, h Handler) {
	e.GET("/page/samples", h.indexViewHandler)
	e.GET("/api/samples", h.apiListSampleHandler)
	e.POST("/api/samples", h.apiNewSampleHandler)
	e.DELETE("/api/samples/:id", h.apiDeleteSampleHandler)

	e.GET("/", h.listPost)
	e.GET("/page/posts", h.listPost)
	e.GET("/api/posts", h.listPost)
	e.POST("/api/posts", h.createPost)
	e.DELETE("/api/posts/:id", h.deletePost)
}
