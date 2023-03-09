package handlers

import (
	"log"
	"net/http"
)

// MainHandler - обработчик всех запросов
func MainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		log.Printf("Поступил POST-запрос")
		ShortenURL(w, r)
	case http.MethodGet:
		log.Println("Поступил GET-запрос")
		GetFullURL(w, r)
	default:
		http.Error(w, "Only GET & POST methods are allowed!", http.StatusBadRequest)
		log.Printf("Поступил некорректный метод запроса: %s\n", r.Method)
	}
}
