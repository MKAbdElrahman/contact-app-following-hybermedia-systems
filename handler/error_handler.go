package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.Logger().Error(err)
	errorPage := fmt.Sprintf("view/error_pages/HTTP%d.html", code)
	fmt.Println(errorPage)
	if err := c.File(errorPage); err != nil {
		c.Logger().Error(err)
	}
}
