package handler

import (
	"app/domain"
	"app/service"
	"app/view"
	"fmt"
	"net/http"

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

func (h *contactHandler) HandleGetAddContact(c echo.Context) error {
	return h.view.RenderAddContactPage(c, view.AddContactPageData{})
}

func (h *contactHandler) HandlePostAddContact(c echo.Context) error {

	var contact domain.Contact

	contact.FirstName = c.FormValue("firstName")
	contact.LastName = c.FormValue("lastName")
	contact.Email = c.FormValue("email")
	contact.Phone = c.FormValue("phone")

	err := h.contactService.ContactStore.AddContact(c.Request().Context(), contact)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return c.Redirect(http.StatusFound, "/contacts")
}
