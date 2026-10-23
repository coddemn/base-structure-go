# === Этап 1: сборка ===
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Кэшируем зависимости: копируем только mod-файлы
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# ГАРАНТИРОВАННО удаляем старую папку docs, чтобы избежать конфликтов
RUN rm -rf docs

# Устанавливаем swag и генерируем docs.go ПЕРЕД сборкой
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12
RUN swag init -g cmd/app/main.go -o docs

# Собираем статические бинарники
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /app/bin/app \
    ./cmd/app

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /app/bin/migrate \
    ./cmd/migrate

# === Этап 2: финальный образ для API ===
FROM alpine:3.20 AS api-image

RUN apk --no-cache add ca-certificates
ENV LANG=C.UTF-8
ENV LC_ALL=C.UTF-8

WORKDIR /app

COPY --from=builder /app/bin/app /app/bin/app
COPY --from=builder /app/configs /app/configs

# === Этап 3: финальный образ для миграций ===
FROM alpine:3.20 AS migrate-image

RUN apk --no-cache add ca-certificates
ENV LANG=C.UTF-8
ENV LC_ALL=C.UTF-8

WORKDIR /app

COPY --from=builder /app/bin/migrate /app/bin/migrate
COPY --from=builder /app/configs /app/configs
