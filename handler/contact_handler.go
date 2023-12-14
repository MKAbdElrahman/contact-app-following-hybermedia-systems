package handler

import (
	"app/domain"
	"app/service"
	"app/view"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

func (h *contactHandler) HandleGetContactByID(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return err
	}

	data, err := h.contactService.ContactStore.GetContactByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrorContactNotFound) {
			return echo.ErrNotFound
		} else {
			return echo.ErrInternalServerError
		}

	}

	return h.view.RenderViewContactPage(c, view.ViewContactPageData{
		Contact: *data,
	})
}

func (h *contactHandler) ValidateEmail(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return echo.ErrBadRequest
	}

	newEmail := strings.TrimSpace(c.FormValue("email"))
	fmt.Println(newEmail)

	contact, _ := h.contactService.ContactStore.GetContactByEmail(c.Request().Context(), newEmail)

	fmt.Println(contact)

	if contact != nil { // there is a contact found with this email
		if contact.ID != id {
			return h.view.RenderValidationError(c, "email already taken")
		}
	}

	return nil
}

func (h *contactHandler) HandleGetAddContact(c echo.Context) error {
	return h.view.RenderAddContactPage(c, view.AddContactPageData{})
}

func (h *contactHandler) HandleGetEditPage(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return err
	}

	contact, err := h.contactService.ContactStore.GetContactByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrorContactNotFound) {
			return echo.ErrNotFound
		} else {
			return echo.ErrInternalServerError
		}

	}

	return h.view.RenderEditContactPage(c, view.EditContactPageData{

		Contact: *contact,
	})
}

func (h *contactHandler) HandlePostedContactEdit(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return err
	}

	var contact service.ContactUpdateParams

	contact.FirstName = c.FormValue("firstName")
	contact.LastName = c.FormValue("lastName")
	contact.Email = c.FormValue("email")
	contact.Phone = c.FormValue("phone")

	err = h.contactService.ContactStore.UpdateContact(c.Request().Context(), id, contact)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return c.Redirect(http.StatusFound, "/contacts")
}

func (h *contactHandler) HandlePostedContactDelete(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return err
	}

	err = h.contactService.ContactStore.DeleteContact(c.Request().Context(), id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/contacts")
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

func (h *contactHandler) HandleGetSearchContactsPage(c echo.Context) error {
	return h.view.RenderSearchContactsPage(c, view.SearchContactsPageData{})
}

func (h *contactHandler) HandlePostSearchContactsPage(c echo.Context) error {

	data, err := h.contactService.ContactStore.GetContacts(c.Request().Context(), c.FormValue("q"))
	if err != nil {
		return err
	}

	return h.view.RenderContactsPage(c, view.ContactsPageData{
		Contacts: data,
		Query:    c.QueryParam("q"),
	})
}
