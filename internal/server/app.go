package server

import (
	"context"
	"github.com/FortovEgor/url-shortener/internal/configs"
	"github.com/FortovEgor/url-shortener/internal/handlers"
	"github.com/FortovEgor/url-shortener/internal/storage"
	"github.com/caarlos0/env/v6"
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

	var cfg configs.Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	//log.Println("PORT:", cfg.Port, "|", cfg.ServerAddress)

	//////////////////////////////////////////////////////////////
	db := storage.NewDatabase()
	h := handlers.NewHandler(db, cfg)
	r.Route("/", func(r chi.Router) {
		r.Get("/{shortURL}", h.GetFullURL)
		r.Post("/api/shorten", h.ShortenJSONURL)
		r.Post("/", h.ShortenURL)
	})
	//////////////////////////////////////////////////////////////

	//port := cfg.Port
	//adr := cfg.ServerAddress
	//parts := strings.Split(adr, ":")[2]
	////temp := string(parts[len(parts)-1])
	//port := strings.Split(parts, "/")[0]
	//fmt.Println("PORT!!!!:", port)

	server := &http.Server{
		Addr:           cfg.ServerAddress, // единственное место в нашем сервере, где используется Порт
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
