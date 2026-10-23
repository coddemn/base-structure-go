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

# RUN go build -o /bin/app ./cmd/app
RUN go build -o /bin/migrate ./cmd/migrate

# Собираем статический бинарник
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /app/bin/app \
    ./cmd/app

# === Этап 2: финальный образ ===
FROM alpine:3.20

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/bin/app /app/bin/app
COPY --from=builder /bin/migrate /bin/migrate
COPY --from=builder /app/configs /app/configs
#COPY --from=builder /app/docs /app/docs
