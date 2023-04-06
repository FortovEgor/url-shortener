package persistent

import (
	"bufio"
	"github.com/FortovEgor/url-shortener/internal/storage"
	"log"
	"os"
	"strings"
)

// Persistent - структура "двойной" БД:
// 1-ая лежит в файле на диске ПК (i.e. persistent)
// 2-ая лежит в ОЗУ ПК (когда сервер запущен)
// Во 2-ую БД происходит копирование данных при старте сервера
// Во 2-ую БД происходит запись данных при обработке очередного запроса
// на добавление short_url
// ПРИМЕЧАНИЕ: файл БД (1-ая БД) имеет несколько строк вида "full_url short_url"
// (т.е. разделитель двух значений - пробел)
type Persistent struct {
	PathToFileStorage string
	VirtualDatabase   *storage.Database
}

func NewStorage(storagePath string) *Persistent {
	database := storage.NewDatabase() // для VirtualDatabase

	if err := loadURLsFromFile(database, storagePath); err != nil {
		log.Fatal(err)
	}

	return &Persistent{
		PathToFileStorage: storagePath,
		VirtualDatabase:   database,
	}
}

// loadURLsFromFile - ф-ия загружает данные из Локальной БД (т.е. файла) в runtime БД (ОЗУ)
func loadURLsFromFile(database *storage.Database, storagePath string) (err error) {
	if storagePath == "" {
		return
	}

	file, err := os.OpenFile(storagePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Err() != nil {
			return
		}
		line := scanner.Text()
		var data []string
		data = strings.Split(line, " ")

		database.AddItemBothValuesKnown(data[0], data[1])
	}
	return
}

func (s *Persistent) GetItem(fullURL string) (string, error) {
	value, err := s.VirtualDatabase.GetItem(fullURL)

	if err != nil {
		return "", err
	}

	return value, nil
}

func (s *Persistent) AddItem(key string) error {
	s.VirtualDatabase.AddItem(key)
	value := storage.MakeShortURLFromFullURL(key)

	if s.PathToFileStorage != "" {
		file, err := os.OpenFile(s.PathToFileStorage, os.O_WRONLY|os.O_APPEND, 0777)
		if err != nil {
			return err
		}

		defer func() {
			_ = file.Close()
		}()

		writer := bufio.NewWriter(file)
		defer func() {
			_ = writer.Flush()
		}()

		// записываем очередную запись в файл
		if _, err = writer.WriteString(key + " " + value + "\n"); err != nil {
			return err
		}
	}
	return nil
}
