package handlers

import (
	"log"
	"net/http"
)

// GetFullURL - Обработчик GET запросов - запросов на получение full_URL из short_URL
func (h *Handler) GetFullURL(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Accept-Encoding", "gzip")
	r.Header.Set("Content-Encoding", "gzip")
	log.Println("path:", r.URL.Path)
	param := r.URL.Path[1:]
	//param := chi.URLParam(r, "id")
	log.Println("param: ", param)

	if param == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Введите идентификатор URL!"))
		if err != nil {
			log.Fatal("Ошибка в записи ответа!")
			return
		}
		return
	}

	target, err := h.db.GetItem(param)
	log.Println("target:", target)
	//storage.URLDB.GetItem(param)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Такого short_url нет в БД!"))
		if err != nil {
			log.Fatal("Ошибка при записи ответа!")
			return
		}
		return
	}

	w.Header().Set("Location", target)
	log.Println(target)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
