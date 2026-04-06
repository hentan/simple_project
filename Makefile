SHELL := /bin/sh
.DEFAULT_GOAL := help

APP_SERVICE := app
DB_SERVICE := postgres

BACKEND_DIR := backend
FRONTEND_DIR := frontend

ENV_FILE ?= $(BACKEND_DIR)/.env

COMPOSE_BE := docker compose --env-file $(ENV_FILE) -f $(BACKEND_DIR)/docker-compose.yaml
COMPOSE_FE := docker compose -f $(FRONTEND_DIR)/docker-compose.yaml

GOOSE_VERSION ?= v3.21.1
GOOSE_BIN := $(CURDIR)/$(BACKEND_DIR)/bin/goose
MIGRATIONS_DIR := $(CURDIR)/$(BACKEND_DIR)/internal/migrations

# Загружаем переменные из backend/.env в shell-команды.
# Удаляем '\r' на лету, чтобы .env с Windows-окончаниями (CRLF) не ломал /bin/sh.
RUN_ENV = \
	tmp="$$(mktemp)"; \
	trap 'rm -f "$$tmp"' EXIT INT TERM; \
	tr -d '\r' < "$(ENV_FILE)" > "$$tmp"; \
	set -a; . "$$tmp"; set +a;

.PHONY: help \
	up up-be up-fe \
	down down-be down-fe \
	restart restart-be restart-fe \
	logs logs-db sh psql \
	up-db up-app wait-db bootstrap-be \
	check-env goose-install goose-version \
	migrate-up migrate-down migrate-status migrate-version \
	migrate-up-one migrate-up-to migrate-down-to migrate-redo migrate-reset \
	migrate-create

help:
	@echo "up                - поднять backend и frontend"
	@echo "up-be             - поднять postgres, накатить миграции, поднять app"
	@echo "up-fe             - поднять frontend"
	@echo "down              - остановить backend и frontend"
	@echo "down-be           - остановить backend и удалить volume"
	@echo "down-fe           - остановить frontend и удалить volume"
	@echo "restart           - перезапустить backend и frontend"
	@echo "restart-be        - перезапустить backend без удаления volume"
	@echo "restart-fe        - перезапустить frontend"
	@echo "logs              - логи app"
	@echo "logs-db           - логи postgres"
	@echo "sh                - shell в app"
	@echo "psql              - psql в postgres"
	@echo "goose-install     - установить goose локально"
	@echo "goose-version     - версия локального goose"
	@echo "migrate-up        - накатить все pending миграции"
	@echo "migrate-down      - откатить последнюю миграцию"
	@echo "migrate-status    - статус миграций"
	@echo "migrate-version   - текущая версия БД"
	@echo "migrate-up-one    - накатить одну миграцию"
	@echo "migrate-up-to     - накатить до версии: make migrate-up-to version=20260315145337"
	@echo "migrate-down-to   - откатить до версии: make migrate-down-to version=20260315145337"
	@echo "migrate-redo      - переиграть последнюю миграцию"
	@echo "migrate-reset     - откатить все миграции"
	@echo "migrate-create    - создать миграцию: make migrate-create name=create_users_table"

up: up-be up-fe

# Важно: backend поднимаем через bootstrap,
# чтобы приложение стартовало после миграций.
up-be: bootstrap-be

up-fe:
	$(COMPOSE_FE) up --build -d

down: down-fe down-be

down-be:
	$(COMPOSE_BE) down

down-fe:
	$(COMPOSE_FE) down -v

restart: restart-be restart-fe

restart-be:
	$(COMPOSE_BE) down
	$(MAKE) up-be

restart-fe:
	$(COMPOSE_FE) down
	$(COMPOSE_FE) up --build -d

logs:
	$(COMPOSE_BE) logs -f $(APP_SERVICE)

logs-db:
	$(COMPOSE_BE) logs -f $(DB_SERVICE)

sh:
	$(COMPOSE_BE) exec $(APP_SERVICE) sh

psql:
	$(COMPOSE_BE) exec $(DB_SERVICE) sh -c 'psql -U "$$POSTGRES_USER" -d "$$POSTGRES_DB"'

up-db:
	$(COMPOSE_BE) up -d $(DB_SERVICE)

up-app:
	$(COMPOSE_BE) up -d $(APP_SERVICE)

wait-db:
	@until $(COMPOSE_BE) exec -T $(DB_SERVICE) sh -c 'pg_isready -U "$$POSTGRES_USER" -d "$$POSTGRES_DB"' >/dev/null 2>&1; do \
		echo "waiting for postgres..."; \
		sleep 1; \
	done; \
	echo "postgres is ready"

bootstrap-be: up-db wait-db migrate-up up-app

check-env:
	@test -f "$(ENV_FILE)" || (echo "$(ENV_FILE) not found"; exit 1)
	@$(RUN_ENV) \
	test -n "$$POSTGRES_USER" || (echo "POSTGRES_USER is not set in $(ENV_FILE)"; exit 1); \
	test -n "$$POSTGRES_PASSWORD" || (echo "POSTGRES_PASSWORD is not set in $(ENV_FILE)"; exit 1); \
	test -n "$$POSTGRES_DB" || (echo "POSTGRES_DB is not set in $(ENV_FILE)"; exit 1); \
	test -n "$$POSTGRES_PORT" || (echo "POSTGRES_PORT is not set in $(ENV_FILE)"; exit 1); \
	test -n "$$DATABASE_DSN" || (echo "DATABASE_DSN is not set in $(ENV_FILE)"; exit 1); \
	test -n "$$DATABASE_DSN_LOCAL" || (echo "DATABASE_DSN_LOCAL is not set in $(ENV_FILE)"; exit 1)

$(GOOSE_BIN):
	@mkdir -p "$(dir $(GOOSE_BIN))"
	@GOBIN="$(dir $(GOOSE_BIN))" go install github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION)

goose-install: $(GOOSE_BIN)

goose-version: $(GOOSE_BIN)
	@"$(GOOSE_BIN)" -version

migrate-up: $(GOOSE_BIN) check-env
	@$(RUN_ENV) \
	GOOSE_DRIVER=postgres \
	GOOSE_DBSTRING="$$DATABASE_DSN_LOCAL" \
	GOOSE_MIGRATION_DIR="$(MIGRATIONS_DIR)" \
	"$(GOOSE_BIN)" up

migrate-down: $(GOOSE_BIN) check-env
	@$(RUN_ENV) \
	GOOSE_DRIVER=postgres \
	GOOSE_DBSTRING="$$DATABASE_DSN_LOCAL" \
	GOOSE_MIGRATION_DIR="$(MIGRATIONS_DIR)" \
	"$(GOOSE_BIN)" down

migrate-status: $(GOOSE_BIN) check-env
	@$(RUN_ENV) \
	GOOSE_DRIVER=postgres \
	GOOSE_DBSTRING="$$DATABASE_DSN_LOCAL" \
	GOOSE_MIGRATION_DIR="$(MIGRATIONS_DIR)" \
	"$(GOOSE_BIN)" status

migrate-version: $(GOOSE_BIN) check-env
	@$(RUN_ENV) \
	GOOSE_DRIVER=postgres \
	GOOSE_DBSTRING="$$DATABASE_DSN_LOCAL" \
	GOOSE_MIGRATION_DIR="$(MIGRATIONS_DIR)" \
	"$(GOOSE_BIN)" version

migrate-up-one: $(GOOSE_BIN) check-env
	@$(RUN_ENV) \
	GOOSE_DRIVER=postgres \
	GOOSE_DBSTRING="$$DATABASE_DSN_LOCAL" \
	GOOSE_MIGRATION_DIR="$(MIGRATIONS_DIR)" \
	"$(GOOSE_BIN)" up-by-one

migrate-up-to: $(GOOSE_BIN) check-env
	@test -n "$(version)" || (echo 'usage: make migrate-up-to version=20260315145337'; exit 1)
	@$(RUN_ENV) \
	GOOSE_DRIVER=postgres \
	GOOSE_DBSTRING="$$DATABASE_DSN_LOCAL" \
	GOOSE_MIGRATION_DIR="$(MIGRATIONS_DIR)" \
	"$(GOOSE_BIN)" up-to "$(version)"

migrate-down-to: $(GOOSE_BIN) check-env
	@test -n "$(version)" || (echo 'usage: make migrate-down-to version=20260315145337'; exit 1)
	@$(RUN_ENV) \
	GOOSE_DRIVER=postgres \
	GOOSE_DBSTRING="$$DATABASE_DSN_LOCAL" \
	GOOSE_MIGRATION_DIR="$(MIGRATIONS_DIR)" \
	"$(GOOSE_BIN)" down-to "$(version)"

migrate-redo: $(GOOSE_BIN) check-env
	@$(RUN_ENV) \
	GOOSE_DRIVER=postgres \
	GOOSE_DBSTRING="$$DATABASE_DSN_LOCAL" \
	GOOSE_MIGRATION_DIR="$(MIGRATIONS_DIR)" \
	"$(GOOSE_BIN)" redo

migrate-reset: $(GOOSE_BIN) check-env
	@$(RUN_ENV) \
	GOOSE_DRIVER=postgres \
	GOOSE_DBSTRING="$$DATABASE_DSN_LOCAL" \
	GOOSE_MIGRATION_DIR="$(MIGRATIONS_DIR)" \
	"$(GOOSE_BIN)" reset

migrate-create: $(GOOSE_BIN)
	@test -n "$(name)" || (echo 'usage: make migrate-create name=create_users_table'; exit 1)
	@"$(GOOSE_BIN)" -dir "$(MIGRATIONS_DIR)" create "$(name)" sql