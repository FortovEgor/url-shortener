package server

import (
	"context"
	"github.com/FortovEgor/url-shortener/internal/configs"
	"github.com/FortovEgor/url-shortener/internal/handlers"
	"github.com/FortovEgor/url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func StartServer() {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	//////////////////////////////////////////////////////////////
	db := storage.NewDatabase()
	h := handlers.NewHandler(db)
	r.Route("/", func(r chi.Router) {
		r.Get("/{shortURL}", h.GetFullURL)
		r.Post("/", h.ShortenURL)
	})
	//////////////////////////////////////////////////////////////

	server := &http.Server{
		Addr:           configs.Port,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := server.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}
