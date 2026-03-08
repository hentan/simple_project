APP_SERVICE=app
DB_SERVICE=postgres
GOOSE=goose

COMPOSE_BE=docker compose -f backend/docker-compose.yaml
COMPOSE_FE=docker compose -f frontend/docker-compose.yaml

.PHONY: up up-be up-fe down down-be down-fe

up: up-be up-fe
up-be:
	$(COMPOSE_BE) up -d
up-fe:
	$(COMPOSE_FE) up --build -d

down: down-fe down-be
down-be:
	$(COMPOSE_BE) down -v
down-fe:
	$(COMPOSE_FE) down -v

restart: restart-be restart-fe
restart-be:
	$(COMPOSE_BE) down
	$(COMPOSE_BE) up -d
restart-fe:
	$(COMPOSE_FE) down
	$(COMPOSE_FE) up -d

logs:
	$(COMPOSE_BE) logs -f $(APP_SERVICE)

logs-db:
	$(COMPOSE_BE) logs -f $(DB_SERVICE)

sh:
	$(COMPOSE_BE) exec $(APP_SERVICE) sh

psql:
	$(COMPOSE_BE) exec $(DB_SERVICE) sh -c 'psql -U "$$POSTGRES_USER" -d "$$POSTGRES_DB"'

migrate-up:
	$(COMPOSE_BE) exec $(APP_SERVICE) sh -c '$(GOOSE) -dir ./migrations postgres "$$DATABASE_DSN" up'

migrate-down:
	$(COMPOSE_BE) exec $(APP_SERVICE) sh -c '$(GOOSE) -dir ./migrations postgres "$$DATABASE_DSN" down'

migrate-status:
	$(COMPOSE_BE) exec $(APP_SERVICE) sh -c '$(GOOSE) -dir ./migrations postgres "$$DATABASE_DSN" status'

