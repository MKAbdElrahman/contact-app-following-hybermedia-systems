package main

import (
	"app/db"
	"app/handler"
	"app/service"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ardanlabs/conf/v3"
)

func main() {

	// -------------------------------------------------------------------------
	// Configuration

	var cfg handler.ServerConfig

	help, err := conf.Parse("", &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return
		}
		log.Fatal("parsing config: %w", err)
	}

	contactStore := db.NewInMemoryContactStore()

	contactStore.SeedMe(context.Background())

	contactService := service.NewContactService(contactStore)

	services := handler.Services{
		ContactService: contactService,
	}
	e := handler.NewServer(context.Background(), cfg, services)

	// Start server
	go func() {
		if err := e.StartServer(e.Server); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
