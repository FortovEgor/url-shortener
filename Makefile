.PHONY: build run test clean all

CC = go

build:  # Сборка проекта под текущую ОС
	$(CC) build -o bin/main cmd/shortener/main.go

run:  # Запускаем сервер
	$(CC) run cmd/shortener/main.go

test:  # Запускаем локальные тесты
	cd internal/handlers && go test

clean:
	cd bin && rm *

all: build
