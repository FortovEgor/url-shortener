package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

// ShortenURL - Обработчик POST запросов - запросов на сокращение URL и занесение в БД сайта
func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Accept-Encoding", "gzip")
	r.Header.Set("Content-Encoding", "gzip")
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

	log.Println("ADDING ITEM - START") // +
	shortURL, _ := h.db.AddItem(param)
	//storage.URLDB.AddItem(param)
	log.Println("ADDING ITEM - END") // -

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(fmt.Sprintf("%s/%s", h.conf.BaseURL, shortURL)))
	if err != nil {
		log.Fatal("Ошибка при записе ответа!")
	}
}
