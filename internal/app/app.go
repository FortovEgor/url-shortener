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

// var URLDB = make(map[uint32]string) // словарь типа "идентификатор_сокращенного_URL : полный_URL"
// сокращенный URL - это host + unique_id(uint32)
var URLDB = make(map[string]string) // словарь типа "short_url:full_url"

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
//func generateIDOfURL(url string) uint32 {
//	return hash(url)
//}

// PerformSeedingOfDB - Ф-ия, заполняющая нашу БД произвольными данными ДО запуска роутераы
func PerformSeedingOfDB() {
	//URLDB[generateIDOfURL("short_url")] = "google.com"
	URLDB["short_url"] = "google.com"
	//log.Println("short-url:", generateIDOfURL("short_url")) // 1806399902
}

func shortenURL(fullURL string) string {
	return "temp.com"
	//var shortURL = Host
	//var uniqueNumber uint32 = 0
	//var val string
	//for _, ok := URLDB[uniqueNumber]; ok; uniqueNumber++ {
	//	val, ok = URLDB[uniqueNumber]
	//	if ok && val == fullURL {
	//		return shortURL + "/" + fmt.Sprint(uniqueNumber)
	//	}
	//	//fmt.Println(uniqueNumber)
	//}
	//if uniqueNumber != 0 {
	//	uniqueNumber--
	//}
	//URLDB[uniqueNumber] = fullURL // добавляем новую запись в нашу БД
	////log.Println("uniqueNumber:", uniqueNumber)
	//return shortURL + "/" + fmt.Sprint(uniqueNumber)
}

// MainHandler - обработчик GET запросов
func MainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		log.Println("Поступил GET-запрос")

		log.Println("path:", r.URL.Path)
		//param := strings.TrimPrefix(r.URL.Path, "/")
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
		http.Error(w, "Only GET & POST methods are allowed!", http.StatusMethodNotAllowed)
		log.Printf("Поступил некорректный метод запроса: %s\n", r.Method)
		return
	}
}
