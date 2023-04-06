package handlers

import "github.com/FortovEgor/url-shortener/internal/configs"

type DatabaseInterface interface {
	GetItem(itemID string) (string, error)
	AddItem(fullURL string) (string, error)
}

// NewHandler - универсальный хендлер для разных БД
func NewHandler(db DatabaseInterface, cfg configs.Config) *Handler {
	return &Handler{db: db, conf: cfg}
}

type Handler struct {
	db   DatabaseInterface
	conf configs.Config
}
