package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

/*
Напишите сервис для сокращения длинных URL. Требования:
- Сервер должен быть доступен по адресу: http://localhost:8080. +
- Сервер должен предоставлять два эндпоинта: POST / и GET /{id}.
- Эндпоинт POST / принимает в теле запроса строку URL для сокращения и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
- Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location. +
- Нужно учесть некорректные запросы и возвращать для них ответ с кодом 400.
*/

const host = "localhost:8080"

var UrlDb = make(map[uint32]string) // словарь типа "идентификатор_сокращенного_URL : полный_URL"

// Jenkins hash function
// source: https://dev.to/ishankhare07/you-think-you-understand-key-value-pairs-m7l
func hash(key string) (hash uint32) {
	hash = 0
	for _, ch := range key {
		hash += uint32(ch)
		hash += hash << 10
		hash ^= hash >> 6
	}

	hash += hash << 3
	hash ^= hash >> 11
	hash += hash << 15

	return
}

// Ф-ия возвращает уникальный идентификатор ресурса
func generateIdOfURL(url string) uint32 {
	return hash(url)
}

// Ф-ия, заполняющая нашу БД произвольными данными ДО запуска роутераы
func performSeedingOfDB() {
	UrlDb[generateIdOfURL("short_url")] = "google.com"
}

// MainHandler - обработчик GET запросов
func MainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		log.Printf("Поступил GET-запрос")
		param := strings.TrimPrefix(r.URL.Path, "/")

		if param == "" {
			http.Error(w, "Введите идентификатор URL!", http.StatusBadRequest)
			return
		}

		if id, err := strconv.Atoi(param); err != nil || id < 0 {
			http.Error(w, "Некорректный параметр запроса!", http.StatusBadRequest)
			return
		}
		id, _ := strconv.Atoi(param)
		fullUrl, exists := UrlDb[uint32(id)]
		if exists { // такой идентификатор URL есть в БД
			w.Header().Set("Location", fullUrl)
			w.WriteHeader(307)
			w.Write([]byte(fullUrl))
		} else {
			http.Error(w, "Такого идентификатора URL нет!", http.StatusBadRequest)
			return
		}
	} else if r.Method == http.MethodPost {
		log.Printf("Поступил POST-запрос")
		///////////////////////////////////////////
		// TODO: implement POST method
		// читаем Body
		b, err := io.ReadAll(r.Body)
		// обрабатываем ошибку
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		// продолжаем обработку
		// ...
		resp, err := json.Unmarshal(b)
		if err != nil {
			log.Fatal("Error!")
		}
		log.Print(resp)
		///////////////////////////////////////////
	} else {
		http.Error(w, "Only GET & POST methods are allowed!", http.StatusBadRequest)
		log.Printf("Поступил некорректный метод запроса: %s\n", r.Method)
		return
	}
}

func main() {
	performSeedingOfDB() // сидирование БД

	http.HandleFunc("/", MainHandler)

	server := &http.Server{
		Addr: host,
	}
	log.Fatal(server.ListenAndServe()) // сервер принудительно завершает свою работу
}
