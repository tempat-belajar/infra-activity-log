package main

import (
	"context"
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/akr/infra-activity-log/internal/api"
	"github.com/akr/infra-activity-log/internal/db"
)

//go:embed web/*
var webFiles embed.FS

func mustEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		logger.Error("DB_DSN env var is required")
		os.Exit(1)
	}
	port := mustEnv("HTTP_PORT", "8080")
	uploadDir := mustEnv("UPLOAD_DIR", "/data/uploads")

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		logger.Error("failed to create upload dir", "err", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.New(ctx, dsn)
	if err != nil {
		logger.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}
	defer database.Close()
	logger.Info("connected to database")

	webSub, err := fs.Sub(webFiles, "web")
	if err != nil {
		logger.Error("failed to load embedded web assets", "err", err)
		os.Exit(1)
	}

	h := api.New(database, uploadDir, logger)
	router := api.NewRouter(h, uploadDir, http.FS(webSub))

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("server starting", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Info("shutting down gracefully")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "err", err)
	}
	logger.Info("server stopped")
}
