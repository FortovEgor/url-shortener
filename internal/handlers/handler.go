package handlers

type DatabaseInterface interface {
	GetItem(itemID string) (string, error)
	AddItem(itemID string, value string)
}

//
//func NewHandler(DB DatabaseInterface) *Handler {
//	return &Handler{s: DB}
//}
//
//type Handler struct {
//	s DatabaseInterface
//}
