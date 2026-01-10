APP_SERVICE=app
DB_SERVICE=postgres
COMPOSE=docker compose
GOOSE=goose

.PHONY: build up down restart logs logs-db sh psql migrate-up migrate-down migrate-status

build:
	$(COMPOSE) build

up:
	$(COMPOSE) up -d

down:
	$(COMPOSE) down -v

restart:
	$(COMPOSE) down
	$(COMPOSE) up -d

logs:
	$(COMPOSE) logs -f $(APP_SERVICE)

logs-db:
	$(COMPOSE) logs -f $(DB_SERVICE)

sh:
	$(COMPOSE) exec $(APP_SERVICE) sh

psql:
	$(COMPOSE) exec $(DB_SERVICE) psql -U $$POSTGRES_USER -d $$POSTGRES_DB

# --- migrations (goose) ---
migrate-up:
	$(COMPOSE) exec $(APP_SERVICE) \
	$(GOOSE) -dir ./migrations postgres "$$DATABASE_DSN" up

migrate-down:
	$(COMPOSE) exec $(APP_SERVICE) \
	$(GOOSE) -dir ./migrations postgres "$$DATABASE_DSN" down

migrate-status:
	$(COMPOSE) exec $(APP_SERVICE) \
	$(GOOSE) -dir ./migrations postgres "$$DATABASE_DSN" status
