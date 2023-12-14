package handler

import (
	"app/service"
	"app/view"
	"app/view/blog"
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"
)

type Services struct {
	ContactService *service.ContactService
}

func NewMux(ctx context.Context, logger *slog.Logger, services Services) *echo.Echo {

	e := echo.New()

	view := view.NewView()

	indexPageHandler := IndexPageHandler{}
	contactHandler := NewContactHandler(services.ContactService, view)

	e.Static("/static", "assets")

	e.HTTPErrorHandler = CustomHTTPErrorHandler

	e.GET("/", indexPageHandler.HandleGetIndexPage)
	e.GET("/contacts", contactHandler.HandleGetContacts)

	e.GET("/contacts/search", contactHandler.HandleGetSearchContactsPage)
	e.POST("/contacts/search", contactHandler.HandlePostSearchContactsPage)

	e.GET("/contacts/:id/view", contactHandler.HandleGetContactByID)

	e.GET("/contacts/:id/edit", contactHandler.HandleGetEditPage)
	e.POST("/contacts/:id/edit", contactHandler.HandlePostedContactEdit)

	e.POST("/contacts/:id/delete", contactHandler.HandlePostedContactDelete)

	e.GET("/contacts/new", contactHandler.HandleGetAddContact)
	e.POST("/contacts/new", contactHandler.HandlePostAddContact)

	myBlog := blog.NewBlog()
	blogHandler := NewBlogHandler(myBlog)

	e.GET("/blog", blogHandler.HandleGetIndexPage)
	e.GET("/blog/contact", blogHandler.HandleGetContactPage)

	return e
}
