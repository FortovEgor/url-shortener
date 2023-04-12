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

func NewStorage(storagePath string) (*Persistent, error) {
	//if storagePath == "" {
	//	return nil, errors.New("не указан путь к БД")
	//}

	database := storage.NewDatabase() // для VirtualDatabase

	// загружаем данные из файла ТОЛЬКО если этот файл существует!
	if storagePath != "" {
		if err := loadURLsFromFile(database, storagePath); err != nil {
			log.Fatal(err)
		}
	}

	return &Persistent{
		PathToFileStorage: storagePath,
		VirtualDatabase:   database,
	}, nil
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
		data := strings.Split(line, " ")

		database.AddItemBothValuesKnown(data[0], data[1])
	}
	return
}

func (s *Persistent) GetItem(shortURL string) (string, error) {
	value, err := s.VirtualDatabase.GetItem(shortURL) // достаем зн-ие из БД в ОЗУ

	if err != nil {
		return "", err
	}

	return value, nil
}

// LoadURLsToStorage - ф-ия, кот-ая загружает все записи из БД в файл
func LoadURLsToStorage(s *Persistent, storagePath string) (err error) {
	if s.PathToFileStorage != "" {
		//file, err := os.OpenFile(s.PathToFileStorage, os.O_WRONLY|os.O_APPEND, 0777)
		file, err := os.OpenFile(s.PathToFileStorage, os.O_WRONLY, 0777)
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

		for shortURL, fullURL := range s.VirtualDatabase.URLs {
			// записываем очередную запись в файл
			if _, err = writer.WriteString(fullURL + " " + shortURL + "\n"); err != nil {
				return err
			}
		}

	}
	return
}

func (s *Persistent) AddItem(fullURL string) (string, error) {
	s.VirtualDatabase.AddItem(fullURL)
	shortURL := storage.MakeShortURLFromFullURL(fullURL)

	return shortURL, nil
}
