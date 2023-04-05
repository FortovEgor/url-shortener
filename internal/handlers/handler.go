package handlers

import "github.com/FortovEgor/url-shortener/internal/configs"

type DatabaseInterface interface {
	GetItem(itemID string) (string, error)
	AddItem(fullURL string) string
}

func NewHandler(db DatabaseInterface, cfg configs.Config) *Handler {
	return &Handler{db: db, conf: cfg}
}

type Handler struct {
	db   DatabaseInterface
	conf configs.Config
}
