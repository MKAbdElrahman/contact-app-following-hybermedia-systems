package handler

import (
	"app/service"
	"app/view"

	"github.com/labstack/echo/v4"
)

type contactHandler struct {
	contactService *service.ContactService
	view           *view.View
}

func NewContactHandler(c *service.ContactService, v *view.View) *contactHandler {
	return &contactHandler{
		contactService: c,
		view:           v,
	}
}

func (h *contactHandler) HandleGetContacts(c echo.Context) error {

	data, err := h.contactService.ContactStore.GetContacts(c.Request().Context(), c.QueryParam("q"))
	if err != nil {
		return err
	}

	return h.view.RenderContactsPage(c, view.ContactsPageData{
		Contacts: data,
		Query:    c.QueryParam("q"),
	})
}
