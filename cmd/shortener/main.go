package main

import (
	"github.com/FortovEgor/url-shortener/internal/server"
)

func main() {
	//database := storage.NewDatabase() // вот как это и куда перекинуть??? не понимаю, сидел 2 дня!!!

	server.StartServer()
}
