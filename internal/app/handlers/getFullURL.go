package handlers

import (
	"github.com/FortovEgor/url-shortener/internal/app/storage"
	"log"
	"net/http"
)

// Обработчик GET запросов - запросов на получение full_URL из short_URL
func GetFullURL(w http.ResponseWriter, r *http.Request) {
	log.Println("path:", r.URL.Path)
	param := r.URL.Path[1:]
	log.Println("param: ", param)

	if param == "" {
		http.Error(w, "Введите идентификатор URL!", http.StatusBadRequest)
		return
	}

	target, err := storage.UrlDB.GetItem(param)
	if err != nil {
		http.Error(w, "Такого short_url нет в БД!", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", target)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
