package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"starterkit/internal/database"
	"starterkit/internal/handlers"
	"syscall"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("/etc/.env"); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	conn := database.GetSqlXConnection()
	if conn == nil {
		panic(errors.New("could not connect to database"))
	}

	r := chi.NewRouter()
	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	//Unprotected routes
	r.Group(func(r chi.Router) {
		r.Use(
			middleware.Logger,
		)
	})

	//Protected routes
	r.Group(func(r chi.Router) {
		r.Use(
			middleware.Logger,
		)
		r.Get("/", handlers.IndexHandler)
	})

	killSig := make(chan os.Signal, 1)
	signal.Notify(killSig, os.Interrupt, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    ":2626",
		Handler: r,
	}

	go func() {
		// Start the server with HTTPS (using a self-signed cert here)
		err := srv.ListenAndServeTLS("server.crt", "server.key")
		if errors.Is(err, http.ErrServerClosed) {
			// If the server was closed gracefully, log the info message
			logger.Info("Server shutdown complete")
		} else if err != nil {
			// Log other errors (e.g., if the server fails to start)
			logger.Error("Server error", slog.Any("err", err))
			os.Exit(1)
		}
	}()

	logger.Info("Server started", slog.String("port", os.Getenv("FRONTEND_PORT")))
	<-killSig
	logger.Info("Shutting down server")

	// Create a context with a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", slog.Any("err", err))
		os.Exit(1)
	}

	logger.Info("Server shutdown complete")
}
