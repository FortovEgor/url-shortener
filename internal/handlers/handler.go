package handlers

type DatabaseInterface interface {
	GetItem(itemID string) (string, error)
	AddItem(fullURL string) string
}

func NewHandler(db DatabaseInterface) *Handler {
	return &Handler{db: db}
}

type Handler struct {
	db DatabaseInterface
}

//func (h * Handler) GetFullURL(w http.ResponseWriter, r *http.Request) {
//	h.db.GetItem(param)
//	// ...
//}
//
//func (h * Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
//	h.db..AddItem(param)
//	// ...
//}
