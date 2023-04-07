package server

import (
	"context"
	"flag"
	"fmt"
	"github.com/FortovEgor/url-shortener/internal/configs"
	"github.com/FortovEgor/url-shortener/internal/handlers"
	gzip "github.com/FortovEgor/url-shortener/internal/server/middleware"
	"github.com/FortovEgor/url-shortener/internal/storage/persistent"
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
	r.Use(gzip.GZIPHandler) // ДОБАВИЛИ сжатие трафика

	var cfg configs.Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg)
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "The address where server is deployed")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "Use ths url as prefix to shortened value")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "Save all shortened URLs to the disk")
	flag.Parse()

	fmt.Println("File for DB:", cfg.FileStoragePath)

	//////////////////////////////////////////////////////////////
	//db := storage.NewDatabase()
	//db := persistent.NewStorage(cfg.FileStoragePath)
	fileDB, _ := persistent.NewStorage(cfg.FileStoragePath)

	h := handlers.NewHandler(fileDB, cfg)
	r.Route("/", func(r chi.Router) {
		r.Get("/{shortURL}", h.GetFullURL)
		r.Post("/api/shorten", h.ShortenJSONURL)
		r.Post("/", h.ShortenURL)
	})
	//////////////////////////////////////////////////////////////
	fmt.Println("STORAGE PATH:", cfg.FileStoragePath)
	//////////////////////////////////////////////////////////////

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
