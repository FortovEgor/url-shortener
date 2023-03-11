package server

import (
	"github.com/FortovEgor/url-shortener/internal/app/configs"
	"github.com/FortovEgor/url-shortener/internal/app/handlers"
	"github.com/FortovEgor/url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

/*
ИНКРЕМЕНТ 1:
Напишите сервис для сокращения длинных URL. Требования:
- Сервер должен быть доступен по адресу: http://localhost:8080. +
- Сервер должен предоставлять два эндпоинта: POST / и GET /{id}.
- Эндпоинт POST / принимает в теле запроса строку URL для сокращения и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле. +
- Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location. +
- Нужно учесть некорректные запросы и возвращать для них ответ с кодом 400. +
*/

/*
ИНКРЕМЕНТ 2:
- Покройте сервис юнит-тестами. Сконцентрируйтесь на покрытии тестами эндпоинтов,
чтобы защитить API сервиса от случайных изменений.
*/

/*
ИНКРЕМЕНТ 3:
Используя любой пакет (роутер или фреймворк), совместимый с net/http,
перепишите ваш код.
В данном случае используется пакет chi
*/

// Примечание: при обработке POST запроса НЕ просиходит проверка на существование
// данного ресурса в БД сайта.

//////////////////////////////////////////////////////////////////////

func StartServer() {
	storage.PerformSeedingOfDB() // сидирование БД
	//
	//http.HandleFunc("/", handlers.MainHandler)
	//
	//server := &http.Server{
	//	Addr: configs.Port,
	//}
	//log.Fatal(server.ListenAndServe()) // сервер принудительно завершает свою работу

	r := chi.NewRouter()
	//r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Get("/{shortURL}", handlers.GetFullURL)
		r.Post("/", handlers.ShortenURL)
	})

	log.Fatal(http.ListenAndServe(configs.Port, r))
}
