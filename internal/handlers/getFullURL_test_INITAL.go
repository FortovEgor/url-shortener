package handlers

import (
	"bytes"
	"github.com/FortovEgor/url-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetFullURL(t *testing.T) {
	type args struct {
		URL    string
		method string
		body   string
	}

	type want struct {
		code           string
		locationHeader string
		body           string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Status 307 - request for google.com (test #1)",
			args: args{
				URL:    "/1d5920f4b44b27a802bd77c4f0536f5a",
				method: http.MethodGet,
			},
			want: want{
				code:           "307 Temporary Redirect",
				locationHeader: "google.com",
			},
		},
		{
			name: "Status 307 - request for github.com (test #2)",
			args: args{
				URL:    "/99cd2175108d157588c04758296d1cfc",
				method: http.MethodGet,
			},
			want: want{
				code:           "307 Temporary Redirect",
				locationHeader: "github.com",
			},
		},
		{
			name: "Status 400 - request for not found (test #3)",
			args: args{
				URL:    "/notfound",
				method: http.MethodGet,
			},
			want: want{
				code:           "400 Bad Request",
				locationHeader: "",
				body:           "Такого short_url нет в БД!",
			},
		},
		{
			name: "Status 400 - request for empty (test #4)",
			args: args{
				URL:    "/",
				method: http.MethodGet,
			},
			want: want{
				code:           "400 Bad Request",
				locationHeader: "",
				body:           "Введите идентификатор URL!",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//storage.PerformSeedingOfDB()

			body := new(bytes.Buffer)
			body.WriteString(tt.args.body)
			request := httptest.NewRequest(tt.args.method, tt.args.URL, body)

			w := httptest.NewRecorder()
			//h := http.HandlerFunc(GetFullURL)
			db := storage.NewDatabase()
			hNew := NewHandler(db)
			h := http.HandlerFunc(hNew.GetFullURL)
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.Status, "Wrong Status!")
			assert.Equal(t, tt.want.body, w.Body.String(), "Wrong body!") // тут ошибка
			assert.Equal(t, tt.want.locationHeader, res.Header.Get("Location"), "Unexpected Location header value")
		})
	}
}
