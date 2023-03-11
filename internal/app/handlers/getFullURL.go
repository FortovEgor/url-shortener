package handlers

import (
	"github.com/FortovEgor/url-shortener/internal/app/storage"
	"log"
	"net/http"
)

// GetFullURL - Обработчик GET запросов - запросов на получение full_URL из short_URL
func GetFullURL(w http.ResponseWriter, r *http.Request) {
	log.Println("path:", r.URL.Path)
	param := r.URL.Path[1:]
	log.Println("param: ", param)

	if param == "" {
		//http.Error(w, "Введите идентификатор URL!", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Введите идентификатор URL!"))
		return
	}

	target, err := storage.URLDB.GetItem(param)
	if err != nil {
		//http.Error(w, "Такого short_url нет в БД!", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Такого short_url нет в БД!"))
		return
	}

	w.Header().Set("Location", target)
	log.Println(target)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
