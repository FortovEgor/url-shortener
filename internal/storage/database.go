package storage

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"log"
	"sync"
)

type Database struct {
	URLs map[string]string // словарь типа "short_url:full_url"
	lock sync.RWMutex
}

func NewDatabase() *Database {
	return &Database{
		URLs: make(map[string]string),
	}
}

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
	db.lock.RLock()
	item, found := db.URLs[shortURL]
	if !found {
		db.lock.RUnlock()
		return "", errors.New("такого short url не найдено")
	}
	db.lock.RUnlock()
	return item, nil
}

// AddItem добавляет пару <shortURL: fullURL> в БД
func (db *Database) AddItem(fullURL string) (shortURL string, err error) {
	db.lock.Lock()
	shortURL = MakeShortURLFromFullURL(fullURL)
	db.URLs[shortURL] = fullURL
	log.Println("shortURL:", shortURL)
	db.lock.Unlock()
	return shortURL, nil
}

// AddItem добавляет пару <shortURL: fullURL> в БД, имея ОБА ЗНАЧЕНИЯ
func (db *Database) AddItemBothValuesKnown(fullURL, shortURL string) {
	db.lock.Lock()
	db.URLs[shortURL] = fullURL
	log.Println("shortURL:", shortURL)
	db.lock.Unlock()
}
