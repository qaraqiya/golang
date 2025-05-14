# syntax=docker/dockerfile:1

FROM golang:1.21-alpine

# Установим зависимости
RUN apk add --no-cache git

# Установим рабочую директорию внутри контейнера
WORKDIR /app

# Скопируем go.mod и go.sum и установим зависимости
COPY go.mod ./
RUN go mod download

# Скопируем остальной код
COPY . .

# Перейдём в папку с сервисом и соберём
WORKDIR /app/expense-service

# Соберем Go-приложение
RUN go build -o app .

# Команда по умолчанию — запуск сервера
CMD ["./app"]
