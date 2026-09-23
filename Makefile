.PHONY: help build build-app build-migrate migrate-up migrate-down migrate-status dev docker-up docker-down clean

# === Базовые переменные ===
APP_BIN := bin/app
MIGRATE_BIN := bin/migrate
DOCKER_ENV := .env.docker


# === СПРАВКА ===
help:
	@echo "Доступные цели (с командой make):"
	@echo "  dev              			- запустить локально (миграции + API)"
	@echo "  migrate-up       			- применить миграции локально"
	@echo "  migrate-down     			- откатить миграции локально"
	@echo "  migrate-status   			- статус миграций"
	@echo "  build            			- собрать оба бинарника"
	@echo "  build-app/build-migrate    - собрать один из бинарников"
	@echo "  docker-up        			- поднять всё в Docker (включая миграции)"
	@echo "  docker-down      			- остановить Docker"
	@echo "  clean            			- удалить бинарники"


# === СБОРКА БИНАРНИКОВ ===
build: build-app build-migrate

build-app:
	go build -o $(APP_BIN) ./cmd/app

build-migrate:
	go build -o $(MIGRATE_BIN) ./cmd/migrate


# === ЛОКАЛЬНЫЕ КОМАНДЫ (для разработки на ПК) ===
migrate-up:
	go run ./cmd/migrate -cmd up

migrate-down:
	go run ./cmd/migrate -cmd down

migrate-status:
	go run ./cmd/migrate -cmd status

dev: migrate-up
	go run ./cmd/app


# === DOCKER КОМАНДЫ ===
docker-up:
	docker compose --env-file $(DOCKER_ENV) up --build

docker-down:
	docker compose down


# === ОЧИСТКА (удаление бинарников) ===
clean:
	rm -f $(APP_BIN) $(MIGRATE_BIN)