package main

import (
	"github.com/FortovEgor/url-shortener/internal/app"
	"log"
	"net/http"
)

func main() {
	app.PerformSeedingOfDB() // сидирование БД

	http.HandleFunc("/", app.MainHandler)

	server := &http.Server{
		Addr: app.Port,
	}
	log.Fatal(server.ListenAndServe()) // сервер принудительно завершает свою работу
}
