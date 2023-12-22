package handler

import (
	"app/service"
	"app/view"
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Services struct {
	ContactService *service.ContactService
}

type ServerConfig struct {
	ReadTimeout     time.Duration `conf:"default:5s"`
	WriteTimeout    time.Duration `conf:"default:10s"`
	IdleTimeout     time.Duration `conf:"default:120s"`
	ShutdownTimeout time.Duration `conf:"default:20s"`
	APIHost         string        `conf:"default:0.0.0.0:3000"`
}

func NewServer(ctx context.Context, cfg ServerConfig, services Services) *echo.Echo {

	e := echo.New()
	e.HideBanner = true

	e.Server.ReadTimeout = cfg.ReadTimeout
	e.Server.WriteTimeout = cfg.WriteTimeout
	e.Server.IdleTimeout = cfg.IdleTimeout
	e.Server.Addr = cfg.APIHost

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "time=${time_custom}, method=${method}, uri=${uri}, status=${status}, latency=${latency_human}\n",
		CustomTimeFormat: "15:04:05",
	}))

	view := view.NewView()

	contactHandler := NewContactHandler(services.ContactService, view)

	e.Static("/static", "assets")

	e.HTTPErrorHandler = CustomHTTPErrorHandler

	e.GET("/", contactHandler.HandleGetSearchContactsPage)
	e.GET("/contacts", contactHandler.HandleGetContacts)
	e.GET("/contacts/count", contactHandler.HandleGetContactsCount)

	e.GET("/contacts/search", contactHandler.HandleGetSearchContactsPage)
	e.POST("/contacts/search", contactHandler.HandlePostSearchContactsPage)

	e.GET("/contacts/:id/view", contactHandler.HandleGetContactByID)
	e.GET("/contacts/:id/email", contactHandler.ValidateEmail)

	e.GET("/contacts/:id/edit", contactHandler.HandleGetEditPage)
	e.POST("/contacts/:id/edit", contactHandler.HandlePostedContactEdit)

	e.DELETE("/contacts/:id", contactHandler.HandlePostedContactDelete)

	e.GET("/contacts/new", contactHandler.HandleGetAddContact)
	e.POST("/contacts/new", contactHandler.HandlePostAddContact)

	e.GET("/contacts/archive", contactHandler.HandleGetArchivePage)
	e.POST("/contacts/archive/status", contactHandler.HandlePostArchive)
	e.GET("/contacts/archive/status", contactHandler.HandleGetArchiveStatus)

	return e
}
