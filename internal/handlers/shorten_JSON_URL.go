package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/FortovEgor/url-shortener/internal/configs"
	"github.com/FortovEgor/url-shortener/internal/storage"
	"log"
	"net/http"
)

type JSONRequest struct {
	Url string `json:"url"`
}

type JSONResponse struct {
	Result string `json:"result"`
}

// ShortenJSONURL - обработчик POST запросов, принимающий принимающий
// в теле запроса JSON-объект {"url":"<some_url>"} и возвращающий в ответ
// объект {"result":"<shorten_url>"}.
// ПРИМЕЧАНИЕ: данный метод НЕ использует БД, он просто ВОЗВРАЩАЕТ short_url !
func (h *Handler) ShortenJSONURL(w http.ResponseWriter, r *http.Request) {
	var request JSONRequest

	// не получилось распарсить входящий JSON
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Failed to unparse incoming JSON; ERROR: "+err.Error(), http.StatusBadRequest)
	}

	if request.Url == "" {
		http.Error(w, "URL is empty!", http.StatusBadRequest)
	}

	url := request.Url // вытягиваем значение url из входящего JSON
	log.Println("URL:", url)

	encodedUrl := storage.MakeShortURLFromFullURL(url)
	//fmt.Print("encoded URL:", encodedUrl)
	//fmt.Print("aaa")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := JSONResponse{
		Result: fmt.Sprintf(configs.Host + encodedUrl),
	}
	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, "Error happened during encoding the answer", http.StatusInternalServerError)
		return
	}

}
