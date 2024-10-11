package endpoint

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.Logger().Error(err)

	if c.Render(code, "pages/error.html", displayError{Message: err.Error()}); err != nil {
		c.Logger().Error(err)
	}
}
