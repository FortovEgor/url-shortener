package app

import (
	"fmt"
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
- Эндпоинт POST / принимает в теле запроса строку URL для сокращения и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле. +
- Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location. +
- Нужно учесть некорректные запросы и возвращать для них ответ с кодом 400. +
*/

// Примечание: при обработке POST запроса просиходит проверка на существование
// данного ресурса в БД сайта.

//////////////////////////////////////////////////////////////////////

const Host = "localhost:8080"

var URLDB = make(map[uint32]string) // словарь типа "идентификатор_сокращенного_URL : полный_URL"
// сокращенный URL - это host + unique_id(uint32)

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
func generateIDOfURL(url string) uint32 {
	return hash(url)
}

// PerformSeedingOfDB - Ф-ия, заполняющая нашу БД произвольными данными ДО запуска роутераы
func PerformSeedingOfDB() {
	URLDB[generateIDOfURL("short_url")] = "google.com"
	//log.Println("short-url:", generateIDOfURL("short_url")) // 1806399902
}

func shortenURL(fullURL string) string {
	var shortURL = Host
	var uniqueNumber uint32 = 0
	var val string
	for _, ok := URLDB[uniqueNumber]; ok; uniqueNumber++ {
		val, ok = URLDB[uniqueNumber]
		if ok && val == fullURL {
			return shortURL + "/" + fmt.Sprint(uniqueNumber)
		}
		//fmt.Println(uniqueNumber)
	}
	if uniqueNumber != 0 {
		uniqueNumber--
	}
	URLDB[uniqueNumber] = fullURL // добавляем новую запись в нашу БД
	//log.Println("uniqueNumber:", uniqueNumber)
	return shortURL + "/" + fmt.Sprint(uniqueNumber)
}

// MainHandler - обработчик GET запросов
func MainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		log.Println("Поступил GET-запрос")
		log.Println("path:", r.URL.Path)
		param := strings.TrimPrefix(r.URL.Path, "/")
		log.Println("param: ", param)

		if param == "" {
			http.Error(w, "Введите идентификатор URL!", http.StatusBadRequest)
			return
		}

		if id, err := strconv.Atoi(param); err != nil || id < 0 {
			log.Println("err:", err, "id:", id)
			http.Error(w, "Некорректный параметр запроса!", http.StatusBadRequest)
			return
		}
		id, _ := strconv.Atoi(param)
		fullURL, exists := URLDB[uint32(id)]
		//w.Write([]byte("Hello world" + fullURL))
		if exists { // такой идентификатор URL есть в БД
			w.Header().Set("Location", fullURL) // ВОТ В ЭТОЙ СТРОЧКЕ ПРОБЛЕМА
			log.Println("HERE!")
			w.WriteHeader(307)
			_, err := w.Write([]byte(fullURL))
			if err != nil {
				http.Error(w, "Ошибка во время возвращения ответа!", http.StatusBadRequest)
				return
			}
		} else {
			http.Error(w, "Такого идентификатора URL нет!", http.StatusBadRequest)
			return
		}
	} else if r.Method == http.MethodPost {
		log.Printf("Поступил POST-запрос")
		///////////////////////////////////////////
		// читаем Body
		b, err := io.ReadAll(r.Body) // тип параметра в теле запроса - plain
		// обрабатываем ошибку
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		log.Println(string(b))  // пришедшее значение
		var fullURL = string(b) // full URL, полученный в запросе
		log.Println("Start shortening url")
		shortURL := shortenURL(fullURL)
		log.Println("End shortening url")
		w.WriteHeader(http.StatusCreated) // код ответа - 201
		_, err = w.Write([]byte(shortURL))
		if err != nil {
			http.Error(w, "Ошибка во время записи ответа!", http.StatusBadRequest)
			return
		} // отправляем текстовую строку в теле ответа
		log.Println(shortURL)
		///////////////////////////////////////////
	} else {
		http.Error(w, "Only GET & POST methods are allowed!", http.StatusBadRequest)
		log.Printf("Поступил некорректный метод запроса: %s\n", r.Method)
		return
	}
}
