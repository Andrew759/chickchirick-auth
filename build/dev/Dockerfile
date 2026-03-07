ARG GOLANG_VERSION='1.25.7'

FROM golang:${GOLANG_VERSION}-alpine

# Устанавливаем зависимости и Delve
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh build-base && \
    go install github.com/go-delve/delve/cmd/dlv@master

# Устанавливается рабочая директория
WORKDIR /app

# Копируется go.mod и go.sum из корневой директории
COPY ./go.mod ./go.sum ./

# Загружаются зависимости go
RUN go mod download

# Копируется исходный код
COPY ./ ./

# Компилируем бинарник В КОРЕНЬ (/main), а не в текущую папку (/app/main)
RUN go build -gcflags="all=-N -l" -o /main ./cmd/main.go

EXPOSE 8080 40000

# Запускаем бинарник из корня
CMD ["dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/main", "--continue"]
