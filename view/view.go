package view

import (
	"app/domain"

	"github.com/labstack/echo/v4"
)

type View struct {
}

func NewView() *View {
	return &View{}
}

func (r *View) RenderContactsPage(c echo.Context, data []domain.Contact) error {

	d := ContactsPageData{
		Contacts: data,
	}

	page := Layout("Contacts", ContactsPageBody(d))

	return page.Render(c.Request().Context(), c.Response().Writer)
}
