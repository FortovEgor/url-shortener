package storage

import "errors"

type DatabaseInterface interface {
	GetItem(itemID string) (string, error)
	AddItem(itemID string, value string)
}

type Database struct {
	URLs map[string]string // словарь типа "short_url:full_url"
}

// UrlDB - экземпляр нашей БД
var UrlDB = Database{URLs: map[string]string{}}

// GetItem возвращает full_url по short_url
func (db *Database) GetItem(shortURL string) (string, error) {
	item, found := db.URLs[shortURL]

	if !found {
		return "", errors.New("такого short url не найдено")
	}
	return item, nil
}

// AddItem добавляет пару <shortURL: fullURL> в БД
func (db *Database) AddItem(shortURL string, fullURL string) {
	db.URLs[shortURL] = fullURL
}

// GetSize возвращает количетсво записей в БД
func (db *Database) GetSize() int {
	return len(db.URLs)
}
