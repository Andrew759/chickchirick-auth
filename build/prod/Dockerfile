# Dockerfile References: https://docs.docker.com/engine/reference/builder/

ARG GOLANG_VERSION='1.26.1'

FROM golang:${GOLANG_VERSION}-alpine as builder

#Устанавливаются нужные пакеты
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

# Устанавливается рабочая директория
WORKDIR /app

# Копируется go.mod и go.sum из корневой директории
COPY ./go.mod ./go.sum ./

# Загружаются зависимости go
RUN go mod download

# Копируется исходный код
COPY ./ ./
COPY ./.env .env

# Компилируется приложение
RUN go build -o main ./cmd/main.go

# Создается финальный образ
FROM alpine:latest

# Устанавливается рабочая директория
WORKDIR /root/

# Копируется бинарник из стадии сборки
COPY --from=builder . .

# Открывается порт
EXPOSE 8080

# Запуск бинарника
CMD ["./app/main"]
