include .env
export

export PROJECT_ROOT=$(shell pwd)

env-up: ## Включить target todoapp-postgres docker compose
	@docker compose up -d todoapp-postgres

env-down: ## Отключить target todoapp-postgres docker compose
	@docker compose down todoapp-postgres

env-cleanup: ## Очистить volume файлы окружения
	@read -p "Очистить все volume файлы окружения? Опасность утечки данных. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down todoapp-postgres port-forwarder && \
		rm -rf out/pgdata && \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка окружения отменена"; \
	fi

migrate-create: ## Создать мииграцию
	@if [ -z "$(seq)" ]; then \
		echo "Отсутствует необходимый параметр seq. Пример: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	docker compose run --rm todoapp-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up: ## Поднять миграцию на уровень выше
	@make migrate-action action=up

migrate-down: ## Опустить миграцию на уровень ниже
	@make migrate-action action=down

migrate-action: ## Вспомогательная функция для migrate-up or migrate-down
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует необходимый параметр action. Пример: make migrate-action action=up"; \
		exit 1; \
	fi; \
	docker compose run --rm todoapp-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

env-port-forward: ## Включить target port-forwarder docker compose
	@docker compose up -d port-forwarder

env-port-close: ## Выключить target port-forwarder docker compose
	@docker compose down port-forwarder

todoapp-run: ## Запуск todoapp приложения
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run cmd/todoapp/main.go

help:
	@grep -E '^[a-zA-Z_0-9_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-20s\033[0m\033[37m%s\033[0m\n", $$1, $$2}'

