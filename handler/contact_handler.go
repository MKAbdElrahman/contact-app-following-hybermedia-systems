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
	const pageSize = 10 // Number of contacts per page

	// Retrieve the query parameter and current page from the request

	query := c.QueryParam("q")
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1 // Default to page 1 if the page parameter is not valid
	}

	// Retrieve contacts data with pagination information
	data, err := h.contactService.ContactStore.GetContactsPaginated(c.Request().Context(), query, page, pageSize)
	if err != nil {
		return err
	}

	// Calculate the total number of pages based on the total number of contacts and the page size
	totalContacts, err := h.contactService.ContactStore.GetTotalContacts(c.Request().Context(), query)
	if err != nil {
		return err
	}
	totalPages := (totalContacts + pageSize - 1) / pageSize

	// Render the contacts page with pagination information
	if c.Request().Header.Get("HX-Trigger") != "" {

		fmt.Println(c.Request().Header.Get("HX-Trigger"))
		return h.view.RenderContactsPageWithoutLayout(c, view.ContactsPageData{
			Contacts:    data,
			Query:       query,
			CurrentPage: page,
			TotalPages:  totalPages,
		})
	} else {
		return h.view.RenderContactsPage(c, view.ContactsPageData{
			Contacts:    data,
			Query:       query,
			CurrentPage: page,
			TotalPages:  totalPages,
		})
	}
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

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/contacts/%d/view", id))
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

	return c.Redirect(http.StatusSeeOther, "/contacts/search")
}

func (h *contactHandler) HandlePostAddContact(c echo.Context) error {

	var contact domain.Contact

	contact.FirstName = c.FormValue("firstName")
	contact.LastName = c.FormValue("lastName")
	contact.Email = c.FormValue("email")
	contact.Phone = c.FormValue("phone")

	id, err := h.contactService.ContactStore.AddContact(c.Request().Context(), contact)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/contacts/%d/view", id))
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
