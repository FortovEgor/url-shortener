package storage

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
)

type DatabaseInterface interface {
	GetItem(itemID string) (string, error)
	AddItem(itemID string, value string)
}

type Database struct {
	URLs map[string]string // словарь типа "short_url:full_url"
}

// URLDB - экземпляр нашей БД
var URLDB = Database{URLs: map[string]string{}}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func MakeShortURLFromFullURL(fullURL string) string {
	return GetMD5Hash(fullURL)
}

// PerformSeedingOfDB - Ф-ия, заполняющая нашу БД произвольными данными ДО запуска роутераы
func PerformSeedingOfDB() {
	fullURLs := [3]string{"google.com", "yandex.ru", "github.com"}
	for _, url := range fullURLs {
		URLDB.AddItem(GetMD5Hash(url), url)
	}
}

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
