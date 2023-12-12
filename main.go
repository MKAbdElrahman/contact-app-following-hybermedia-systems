package main

import (
	"app/db"
	"app/domain"
	"app/handler"
	"app/service"
	"app/view"
	"context"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	contactStore := db.NewInMemoryContactStore()
	contacts := []domain.Contact{
		{FirstName: "Mohamed", LastName: "Ali", Email: "mohamed.ali@example.com", Phone: "+211111111111"},
		{FirstName: "Sayed", LastName: "Kamal", Email: "sayed.kamal@example.com", Phone: "+511111111111"},
		{FirstName: "John", LastName: "Doe", Email: "john.doe@example.com", Phone: "+1234567890"},
		{FirstName: "Jane", LastName: "Doe", Email: "jane.doe@example.com", Phone: "+9876543210"},
		{FirstName: "Alice", LastName: "Smith", Email: "alice.smith@example.com", Phone: "+1122334455"},
		{FirstName: "Bob", LastName: "Johnson", Email: "bob.johnson@example.com", Phone: "+9988776655"},
		{FirstName: "Emily", LastName: "Wilson", Email: "emily.wilson@example.com", Phone: "+6677889900"},
		{FirstName: "David", LastName: "Miller", Email: "david.miller@example.com", Phone: "+5544332211"},
		{FirstName: "Sophia", LastName: "Clark", Email: "sophia.clark@example.com", Phone: "+6677001122"},
		{FirstName: "Daniel", LastName: "Jones", Email: "daniel.jones@example.com", Phone: "+1122334455"},
		// Add more contacts as needed...
	}
	for _, contact := range contacts {
		contactStore.AddContact(context.Background(), contact)
	}
	contactService := service.NewContactService(contactStore)

	view := view.NewView()

	indexPageHandler := handler.IndexPageHandler{}
	contactHandler := handler.NewContactHandler(contactService, view)

	e.Static("/static", "assets")

	e.HTTPErrorHandler = handler.CustomHTTPErrorHandler

	e.GET("/", indexPageHandler.HandleGetIndexPage)
	e.GET("/contacts", contactHandler.HandleGetContacts)
	e.GET("/contacts/:id/view", contactHandler.HandleGetContactByID)

	e.GET("/contacts/:id/edit", contactHandler.HandleGetEditPage)
	e.POST("/contacts/:id/edit", contactHandler.HandlePostedContactEdit)

	e.POST("/contacts/:id/delete", contactHandler.HandlePostedContactDelete)

	e.GET("/contacts/new", contactHandler.HandleGetAddContact)
	e.POST("/contacts/new", contactHandler.HandlePostAddContact)

	e.Logger.Fatal(e.Start(":1323"))
}
