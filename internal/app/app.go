package app

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

/*
Напишите сервис для сокращения длинных URL. Требования:
- Сервер должен быть доступен по адресу: http://localhost:8080. +
- Сервер должен предоставлять два эндпоинта: POST / и GET /{id}.
- Эндпоинт POST / принимает в теле запроса строку URL для сокращения и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле. +
- Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location. +
- Нужно учесть некорректные запросы и возвращать для них ответ с кодом 400. +
*/

// Примечание: при обработке POST запроса просиходит проверка на существование
// данного ресурса в БД сайта.

//////////////////////////////////////////////////////////////////////

const Host = "http://localhost:8080/"
const Port = ":8080"

var URLDB = make(map[string]string) // словарь типа "short_url:full_url"

// PerformSeedingOfDB - Ф-ия, заполняющая нашу БД произвольными данными ДО запуска роутераы
func PerformSeedingOfDB() {
	URLDB["short_url"] = "google.com"
}

// MainHandler - обработчик GET запросов
func MainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		log.Println("Поступил GET-запрос")

		log.Println("path:", r.URL.Path)
		param := r.URL.Path[1:]
		log.Println("param: ", param)

		if param == "" {
			http.Error(w, "Введите идентификатор URL!", http.StatusBadRequest)
			return
		}

		target, ok := URLDB[param]
		if !ok {
			http.Error(w, "Такого short_url нет в БД!", http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", target)
		w.WriteHeader(http.StatusTemporaryRedirect)

	} else if r.Method == http.MethodPost {
		log.Printf("Поступил POST-запрос")
		///////////////////////////////////////////
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

		id := fmt.Sprintf("%x", len(URLDB))
		URLDB[id] = param

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(Host + id))
		///////////////////////////////////////////
	} else {
		http.Error(w, "Only GET & POST methods are allowed!", http.StatusBadReques)
		log.Printf("Поступил некорректный метод запроса: %s\n", r.Method)
		return
	}
}
