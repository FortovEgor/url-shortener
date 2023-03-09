package handlers

import (
	"fmt"
	"github.com/FortovEgor/url-shortener/internal/app/server"
	"github.com/FortovEgor/url-shortener/internal/app/storage"
	"io"
	"log"
	"net/http"
	"net/url"
)

// Обработчик POST запросов - запросов на сокращение URL и занесение в БД сайта
func ShortenURL(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body) // тип параметра в теле запроса - plain
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	log.Println(string(b)) // пришедшее значение
	var param = string(b)  // full URL, полученный в запросе

	if _, err := url.Parse(param); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("То, что Вы ввели, непохоже на url!"))
		return
	}

	id := fmt.Sprintf("%x", storage.UrlDB.GetSize())
	storage.UrlDB.AddItem(id, param)

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(server.Host + id))
}
