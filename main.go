package main

import (
	"app/db"
	"app/domain"
	"app/handler"
	"app/service"
	"app/view"

	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleIndex(c echo.Context) error {
	return c.Redirect(http.StatusFound, "/contacts")
}

func main() {
	e := echo.New()

	contactStore := db.NewInMemoryContactStore()
	contactStore.AddContact(context.Background(), domain.Contact{FirstName: "Mohamed"})
	contactStore.AddContact(context.Background(), domain.Contact{FirstName: "Ahmed"})

	contactService := service.NewContactService(contactStore)

	view := view.NewView()

	contactHandler := handler.NewContactHandler(contactService, view)

	e.GET("/", HandleIndex)
	e.GET("/contacts", contactHandler.HandleGetContacts)

	e.Logger.Fatal(e.Start(":1323"))
}
