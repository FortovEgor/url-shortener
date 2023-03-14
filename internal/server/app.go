package server

import (
	"github.com/FortovEgor/url-shortener/internal/configs"
	handlers2 "github.com/FortovEgor/url-shortener/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func StartServer() {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Get("/{shortURL}", handlers2.GetFullURL)
		r.Post("/", handlers2.ShortenURL)
	})

	log.Fatal(http.ListenAndServe(configs.Port, r))
}
