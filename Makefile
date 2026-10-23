SHELL := /bin/bash

.PHONY: help build build-app build-migrate migrate-up migrate-down migrate-status dev docker-up docker-down clean docs test

# === Базовые переменные ===
APP_BIN := bin/app
MIGRATE_BIN := bin/migrate
DOCKER_ENV := .env.docker
MAIN_GO := cmd/app/main.go


# === СПРАВКА ===
help:
	@echo "🛠 Доступные цели (используйте make <цель>):"
	@echo "  dev                - запустить локально (миграции + API)"
	@echo "  migrate-up         - применить миграции локально"
	@echo "  migrate-down       - откатить миграции локально"
	@echo "  migrate-status     - статус миграций"
	@echo "  build              - собрать оба бинарника"
	@echo "  build-app          - собрать только бинарник API"
	@echo "  build-migrate      - собрать только бинарник миграций"
	@echo "  docker-up          - поднять всё в Docker (включая миграции)"
	@echo "  docker-down        - остановить Docker-контейнеры"
	@echo "  clean              - удалить бинарники и папку docs"
	@echo "  docs               - сгенерировать Swagger-документацию"
	@echo "  test               - запустить все тесты"



# === СБОРКА БИНАРНИКОВ ===
build: build-app build-migrate

build-app:
	@echo "🔨 Сборка бинарника API..."
	go build -o $(APP_BIN) ./cmd/app

build-migrate:
	@echo "🔨 Сборка бинарника миграций..."
	go build -o $(MIGRATE_BIN) ./cmd/migrate



# === ЛОКАЛЬНЫЕ КОМАНДЫ (для разработки на ПК) ===
migrate-up:
	@echo "⬆️ Применение миграций..."
	go run ./cmd/migrate -cmd up

migrate-down:
	@echo "⬇️ Откат миграций..."
	go run ./cmd/migrate -cmd down

migrate-status:
	@echo "📊 Статус миграций:"
	go run ./cmd/migrate -cmd status

dev: migrate-up
	@echo "🚀 Запуск приложения в режиме разработки..."
	go run ./cmd/app



# === DOCKER КОМАНДЫ ===
docker-up:
	@echo "🐳 Поднимаем окружение в Docker..."
	docker compose --env-file $(DOCKER_ENV) up --build

docker-down:
	@echo "🛑 Останавливаем Docker-окружение..."
	docker compose down



# === ОЧИСТКА ===
clean:
	@echo "🧹 Очистка артефактов сборки..."
	rm -f $(APP_BIN) $(MIGRATE_BIN)
	rm -rf docs/



# === ДОКУМЕНТАЦИЯ ===
# Требует предустановки: go install github.com/swaggo/swag/cmd/swag@v1.8.12
docs:
	@echo "📝 Генерация Swagger-документации..."
	swag init -g $(MAIN_GO) -o docs
	@echo "✅ Готово! Документация находится в папке docs/"



# === Тесты ===
test:
	@echo "🧪 Запуск тестов..."
	go test -v -cover ./...