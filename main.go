package main

import (
	"app/db"
	"app/handler"
	"app/service"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	err := run(context.Background(), logger)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(ctx context.Context, logger *slog.Logger) error {

	// -------------------------------------------------------------------------
	// Configuration

	cfg := struct {
		conf.Version
		ServerConfig struct {
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s"`
			APIHost         string        `conf:"default:0.0.0.0:3000"`
		}
	}{Version: conf.Version{
		Desc:  "Go+HTMX",
		Build: "0.0.1",
	}}

	help, err := conf.Parse("", &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	contactStore := db.NewInMemoryContactStore()

	contactStore.SeedMe(ctx)

	contactService := service.NewContactService(contactStore)

	services := handler.Services{
		ContactService: contactService,
	}
	e := handler.NewMux(ctx, logger, services)

	// -------------------------------------------------------------------------
	// App Starting

	logger.Info("starting app")
	defer logger.Info("shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	logger.Info("startup", "config", out)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	server := http.Server{
		Addr:         cfg.ServerConfig.APIHost,
		Handler:      e,
		ReadTimeout:  cfg.ServerConfig.ReadTimeout,
		WriteTimeout: cfg.ServerConfig.WriteTimeout,
		IdleTimeout:  cfg.ServerConfig.IdleTimeout,
		ErrorLog:     slog.NewLogLogger(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{}), slog.LevelError),
	}

	serverErrors := make(chan error, 1)

	go func() {
		logger.Info("startup", "status", "app router started", "host", server.Addr)

		serverErrors <- server.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		logger.Info("shutdown", "status", "shutdown started", "signal", sig)
		defer logger.Info("shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, cfg.ServerConfig.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			server.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
