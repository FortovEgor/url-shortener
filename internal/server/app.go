package server

import (
	"github.com/FortovEgor/url-shortener/internal/configs"
	handlers2 "github.com/FortovEgor/url-shortener/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func StartServer() {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Get("/{shortURL}", handlers2.GetFullURL)
		r.Post("/", handlers2.ShortenURL)
	})

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	if err := http.ListenAndServe(configs.Port, r); err != nil {
		log.Fatalf("listen: %s\n", err)
	}
	<-done
	log.Print("Server Stopped")
}
