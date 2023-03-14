build:
	echo "Сборка проекта под текущую ОС..."
	go build -o bin/main cmd/shortener/main.go

run:
	echo "Запускаем сервер..."
	go run cmd/shortener/main.go

test:
	echo "Запускаем локальные тесты..."
	cd internal/handlers && go test

all: build
