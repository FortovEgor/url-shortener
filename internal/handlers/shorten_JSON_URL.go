package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type JSONRequest struct {
	URL string `json:"url"`
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

	if request.URL == "" {
		http.Error(w, "URL is empty!", http.StatusBadRequest)
	}

	url := request.URL // вытягиваем значение url из входящего JSON
	log.Println("URL:", url)

	encodedURL := h.db.AddItem(url)
	//fmt.Print("encoded URL:", encodedURL)
	//fmt.Print("aaa")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := JSONResponse{
		Result: fmt.Sprintf(h.conf.Host + encodedURL),
	}
	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, "Error happened during encoding the answer", http.StatusInternalServerError)
		return
	}

}
