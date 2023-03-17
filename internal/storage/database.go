package storage

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"sync"
)

type Database struct {
	URLs map[string]string // словарь типа "short_url:full_url"
}

func NewDatabase() *Database {
	return &Database{
		URLs: make(map[string]string),
	}
}

// URLDB - экземпляр нашей БД
//var URLDB = NewDatabase()

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func MakeShortURLFromFullURL(fullURL string) string {
	return GetMD5Hash(fullURL)
}

// PerformSeedingOfDB - Ф-ия, заполняющая нашу БД произвольными данными ДО запуска роутераы
//func PerformSeedingOfDB() {
//	fullURLs := [3]string{"google.com", "yandex.ru", "github.com"}
//	for _, url := range fullURLs {
//		URLDB.AddItem(url)
//	}
//}

// GetItem возвращает full_url по short_url
func (db *Database) GetItem(shortURL string) (string, error) {
	lock := sync.RWMutex{}
	lock.RLock()
	item, found := db.URLs[shortURL]
	if !found {
		return "", errors.New("такого short url не найдено")
	}
	lock.RUnlock()
	return item, nil
}

// AddItem добавляет пару <shortURL: fullURL> в БД
func (db *Database) AddItem(fullURL string) (shortURL string) {
	lock := sync.RWMutex{}
	lock.Lock()
	shortURL = MakeShortURLFromFullURL(fullURL)
	db.URLs[shortURL] = fullURL
	lock.Unlock()
	return
}
