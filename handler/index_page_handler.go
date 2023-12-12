package handler

import (
	"app/view"

	"github.com/labstack/echo/v4"
)

type IndexPageHandler struct {
	view *view.View
}

func (h *IndexPageHandler) HandleGetIndexPage(c echo.Context) error {
	return h.view.RenderIndexPage(c, view.IndexPageData{})
}
