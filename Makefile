include .env
export

export PROJECT_ROOT=$(shell pwd)
export POSTGRES_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@todoapp-postgres:5432/$(POSTGRES_DB)?sslmode=disable

.PHONY: help env-up env-down env-cleanup migrate-create migrate-up migrate-down migrate-force migrate-drop migrate-action env-port-db-forward env-port-db-close todoapp-run

help: ## Показать список команд
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z0-9_-]+:.*## / {printf "  %-24s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

env-up: ## Запустить Postgres
	@docker compose up -d todoapp-postgres

env-down: ## Остановить Postgres
	@docker compose down todoapp-postgres

env-cleanup: ## Остановить Postgres и удалить данные, спросив подтверждение
	@read -p "Очисть данные бд? (y/n): " answer; \
	if [ "$$answer" = "y" ]; then \
		docker compose down todoapp-postgres port-forward && \
		rm -rf ${PROJECT_ROOT}/out/pgdata && \
		echo "Данные бд очищены"; \
	else \
		echo "Отмена"; \
	fi

migrate-create: ## Создать файлы новой миграции
	@read -p "Введите название миграции: " name; \
	if [ -z "$$name" ]; then \
		echo "Название миграции не может быть пустым"; \
		exit 1; \
	fi; \
	echo "Создание миграции $$name"; \
	docker compose run --rm todoapp-postgres-migrate \
		create \
		-ext sql \
		-dir migrations \
		-seq "$$name"

migrate-up: ## Применить все новые миграции
	@make migrate-action action=up

migrate-down: ## Откатить последнюю миграцию
	@make migrate-action action=down

migrate-force: ## Записать версию миграции без её выполнения, version=N (по умолчанию 0)
	@echo "force version=$(or $(version),0)"; \
	docker compose run --rm todoapp-postgres-migrate \
		-path migrations \
		-database "$(POSTGRES_URL)" \
		force "$(or $(version),0)"

migrate-drop: ## Удалить все таблицы в базе
	@docker compose run --rm todoapp-postgres-migrate \
		-path migrations \
		-database "$(POSTGRES_URL)" \
		drop -f

migrate-action: ## Выполнить миграцию, action=up или action=down
	@if [ -z "$(action)" ]; then \
		echo "action не может быть пустым"; \
		exit 1; \
	fi; \
	docker compose run --rm todoapp-postgres-migrate \
		-path migrations \
		-database "$(POSTGRES_URL)" \
		"$(action)"

env-port-db-forward: ## Пробросить порт Postgres на localhost
	@docker compose up -d port-forward

env-port-db-close: ## Закрыть проброс порта Postgres
	@docker compose down port-forward

logs-clear: ## Очистить логи приложения
	@read -p "Очистить логи приложения? (y/n): " answer; \
	if [ "$$answer" = "y" ]; then \
		rm -rf $(PROJECT_ROOT)/out/logs/*.log; \
		echo "Логи приложения очищены"; \
	else \
		echo "Отмена"; \
	fi

todoapp-run: ## Обновить зависимости и запустить приложение
	@go mod tidy && \
	go run $(PROJECT_ROOT)/cmd/todoapp/main.go
