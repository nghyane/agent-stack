package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nghiahoang/template-api/internal/platform/config"
	"github.com/nghiahoang/template-api/internal/platform/database"
	"github.com/nghiahoang/template-api/internal/platform/server"
	"github.com/nghiahoang/template-api/internal/platform/session"
	"github.com/nghiahoang/template-api/internal/platform/settings"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "openapi" {
		spec, err := server.OpenAPISpec()
		if err != nil {
			slog.Error("openapi export failed", "error", err)
			os.Exit(1)
		}
		os.Stdout.Write(spec)
		return
	}

	if err := run(); err != nil {
		slog.Error("startup failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx := context.Background()

	if err := database.Migrate(cfg.DatabaseURL); err != nil {
		return err
	}
	slog.Info("migrations applied")

	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	sessions := session.NewManager(pool, cfg.IsProduction())

	settingsMgr, err := settings.NewManager(ctx, pool)
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           server.New(cfg, pool, sessions, settingsMgr),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		slog.Info("server listening", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	slog.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
