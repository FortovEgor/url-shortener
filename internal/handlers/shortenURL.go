package handlers

import (
	"github.com/FortovEgor/url-shortener/internal/configs"
	"github.com/FortovEgor/url-shortener/internal/storage"
	"io"
	"log"
	"net/http"
	"net/url"
)

// ShortenURL - Обработчик POST запросов - запросов на сокращение URL и занесение в БД сайта
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

	shortURL := storage.URLDB.AddItem(param)

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(configs.Host + shortURL))
	if err != nil {
		log.Fatal("Ошибка при записе ответа!")
	}
}
