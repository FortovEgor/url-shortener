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
