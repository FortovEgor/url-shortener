package handlers

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShortenURL(t *testing.T) {
	type args struct {
		URL    string
		method string
		body   string
	}

	type want struct {
		code string
		body string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Status 201 - test for initial values in DB after seeding (test #1)",
			args: args{
				URL:    "/",
				method: http.MethodPost,
				body:   "google.com",
			},
			want: want{
				code: "201 Created",
				body: "http://localhost:8080/1d5920f4b44b27a802bd77c4f0536f5a",
			},
		},
		{
			name: "Status 201 - test for initial values in DB after seeding (test #2)",
			args: args{
				URL:    "/",
				method: http.MethodPost,
				body:   "github.com",
			},
			want: want{
				code: "201 Created",
				body: "http://localhost:8080/99cd2175108d157588c04758296d1cfc",
			},
		},
		{
			name: "Status 201 - test for a new value (test #3)",
			args: args{
				URL:    "/",
				method: http.MethodPost,
				body:   "hse.ru",
			},
			want: want{
				code: "201 Created",
				body: "http://localhost:8080/981596c01241676b727cae531659b7b8",
			},
		},
		{
			name: "Status 201 - test for a new value (test #4)",
			args: args{
				URL:    "/",
				method: http.MethodPost,
				body:   "ru.wikipedia.org/wiki/Динамически_подключаемая_библиотека",
			},
			want: want{
				code: "201 Created",
				body: "http://localhost:8080/b6115a3fb022d7a06288104ed331f0bd",
			},
		},
		{
			name: "Status 400 - test for a new Incorrect value (test #5)",
			args: args{
				URL:    "/",
				method: http.MethodPost,
				body:   "с{\n\t\"aaa\": \"bbb\"\n}",
			},
			want: want{
				code: "400 Bad Request",
				body: "То, что Вы ввели, непохоже на url!",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//storage.PerformSeedingOfDB()
			// отправляем запрос с телом

			body := new(bytes.Buffer)
			body.WriteString(tt.args.body)
			request := httptest.NewRequest(tt.args.method, tt.args.URL, body)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(ShortenURL)
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.Status, "Wrong Status!")
			assert.Equal(t, tt.want.body, w.Body.String(), "Wrong body!")
		})
	}
}
