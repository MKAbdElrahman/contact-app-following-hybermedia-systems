package view

import (
	"github.com/labstack/echo/v4"
)

type View struct {
}

func NewView() *View {
	return &View{}
}

func (r *View) RenderContactsPage(c echo.Context, data ContactsPageData) error {

	page := Layout("Contacts", ContactsPageBody(c.Request().Context(), data))

	return page.Render(c.Request().Context(), c.Response().Writer)
}

func (r *View) RenderAddContactPage(c echo.Context, data AddContactPageData) error {

	page := Layout("New Contact", AddContactPageBody(c.Request().Context(), data))

	return page.Render(c.Request().Context(), c.Response().Writer)
}
