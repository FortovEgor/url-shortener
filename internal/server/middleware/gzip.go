package gzip

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GZIPHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// СЖАТИЕ ДАННЫХ СЕРВЕРОМ
		// проверяем, что клиент поддерживает gzip-сжатие
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// создаём gzip.Writer поверх текущего w
			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				if _, err := io.WriteString(w, err.Error()); err != nil {
					return
				}
				return
			}
			defer func() {
				_ = gz.Close()
			}()

			w.Header().Set("Content-Encoding", "gzip")
			// передаём обработчику страницы переменную типа gzipWriter для вывода данных
			w = gzipWriter{ResponseWriter: w, Writer: gz}
		}

		// РАСПАКОВКА ДАННЫХ СЕРВЕРОМ
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			// создаём *gzip.Reader, который будет читать тело запроса
			// и распаковывать его
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer gz.Close()
			r.Body = gz
		}

		next.ServeHTTP(w, r)
	})
}
