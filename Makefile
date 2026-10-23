include .env
export

export PROJECT_ROOT=$(shell pwd)
export POSTGRES_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@todoapp-postgres:5432/$(POSTGRES_DB)?sslmode=disable

.PHONY: help env-up env-down env-cleanup env-cleanup-postgres env-cleanup-redis migrate-create migrate-up migrate-down migrate-force migrate-drop migrate-action env-port-postgres-forward env-port-postgres-close env-port-redis-forward env-port-redis-close logs-clear todoapp-run todoapp-deploy todoapp-undeploy swagger-gen

help: ## показать список команд
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z0-9_-]+:.*## / {printf "  %-28s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

env-up: ## поднять postgres и redis
	@docker compose up -d todoapp-postgres todoapp-redis

env-down: ## остановить postgres и redis
	@docker compose down todoapp-postgres todoapp-redis

env-cleanup: ## удалить данные postgres и redis
	@read -p "Очисть данные бд? (y/n): " answer; \
	if [ "$$answer" = "y" ]; then \
		docker compose down todoapp-postgres todoapp-redis port-forward-postgres port-forward-redis && \
		rm -rf ${PROJECT_ROOT}/out/pgdata && \
		echo "Данные postgres и redis очищены"; \
	else \
		echo "Отмена"; \
	fi

env-cleanup-postgres: ## удалить данные postgres
	@read -p "Очисть данные бд? (y/n): " answer; \
	if [ "$$answer" = "y" ]; then \
		docker compose down todoapp-postgres port-forward-postgres && \
		rm -rf ${PROJECT_ROOT}/out/pgdata && \
		echo "Данные postgres очищены"; \
	else \
		echo "Отмена"; \
	fi

env-cleanup-redis: ## удалить данные redis
	@read -p "Очисть данные redis? (y/n): " answer; \
	if [ "$$answer" = "y" ]; then \
		docker compose down todoapp-redis port-forward-redis && \
		echo "Данные redis очищены"; \
	else \
		echo "Отмена"; \
	fi

migrate-create: ## создать файл миграции
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

migrate-up: ## применить миграции
	@make migrate-action action=up

migrate-down: ## откатить последнюю миграцию
	@make migrate-action action=down

migrate-force: ## выставить версию миграции без выполнения SQL, version=N
	@echo "force version=$(or $(version),0)"; \
	docker compose run --rm todoapp-postgres-migrate \
		-path migrations \
		-database "$(POSTGRES_URL)" \
		force "$(or $(version),0)"

migrate-drop: ## удалить все таблицы в базе
	@docker compose run --rm todoapp-postgres-migrate \
		-path migrations \
		-database "$(POSTGRES_URL)" \
		drop -f

migrate-action: ## выполнить действие migrate, нужно action=up или action=down
	@if [ -z "$(action)" ]; then \
		echo "action не может быть пустым"; \
		exit 1; \
	fi; \
	docker compose run --rm todoapp-postgres-migrate \
		-path migrations \
		-database "$(POSTGRES_URL)" \
		"$(action)"

env-port-postgres-forward: ## пробросить порт postgres на localhost
	@docker compose up -d port-forward-postgres

env-port-postgres-close: ## закрыть проброс порта postgres
	@docker compose down port-forward-postgres

env-port-redis-forward: ## пробросить порт redis на localhost
	@docker compose up -d port-forward-redis

env-port-redis-close: ## закрыть проброс порта redis
	@docker compose down port-forward-redis

logs-clear: ## удалить логи приложения
	@read -p "Очистить логи приложения? (y/n): " answer; \
	if [ "$$answer" = "y" ]; then \
		rm -rf $(PROJECT_ROOT)/out/logs/*.log; \
		echo "Логи приложения очищены"; \
	else \
		echo "Отмена"; \
	fi

todoapp-run: ## собрать фронтенд и запустить приложение
	@npm --prefix $(PROJECT_ROOT)/web run build && \
	go mod tidy && \
	go run $(PROJECT_ROOT)/cmd/todoapp/main.go

todoapp-deploy: ## собрать и поднять приложение в docker
	@docker compose up -d --build todoapp

todoapp-undeploy: ## остановить приложение в docker
	@docker compose down todoapp

swagger-gen: ## сгенерировать swagger-документацию
	@docker compose run --rm swagger \
		init \
		-g cmd/todoapp/main.go \
		-o docs \
		--parseInternal \
		--parseDependency