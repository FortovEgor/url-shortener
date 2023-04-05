package handlers

import (
	"bytes"
	"github.com/FortovEgor/url-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_ShortenJSONURL(t *testing.T) {
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
		// TODO: Add test cases.
		{
			name: "Simple test (test #1)",
			args: args{
				URL:    "/api/shorten",
				method: http.MethodPost,
				body:   `{"url":"<some_url>"}`,
			},
			want: want{
				code: "201 Created",
				body: "{\"result\":\"http://localhost:8080/956f7bcc47c9639a868bedf5be29fd65\"}\n",
			},
		},
		{
			name: "Simple test (test #2)",
			args: args{
				URL:    "/api/shorten",
				method: http.MethodPost,
				body:   `{"url":"google.com"}`,
			},
			want: want{
				code: "201 Created",
				body: "{\"result\":\"http://localhost:8080/1d5920f4b44b27a802bd77c4f0536f5a\"}\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := new(bytes.Buffer)
			body.WriteString(tt.args.body)
			request := httptest.NewRequest(tt.args.method, tt.args.URL, body)

			w := httptest.NewRecorder()
			db := storage.NewDatabase()
			hNew := NewHandler(db)
			h := http.HandlerFunc(hNew.ShortenJSONURL)
			h.ServeHTTP(w, request)
			res := w.Result()
			//defer func(Body io.ReadCloser) {
			//	err := Body.Close()
			//	if err != nil {
			//		log.Println("Failed to close the body! ERROR: " + err.Error())
			//	}
			//}(res.Body)
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.Status, "Wrong Status!")
			assert.Equal(t, tt.want.body, w.Body.String(), "Wrong body!")
		})
	}
}
