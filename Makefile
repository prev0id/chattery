API_PATH=./api
MIGRATIONS_PATH=migrations
POSTGRES_STRING=postgresql://user:password@localhost:5432/chattery?sslmode=disable
DOCKER_COMPOSE_BIN=docker-compose

.PHONY: run
run:
	go tool air -c .air.toml

.PHONY: build
build:
	go build -o ./bin/chattery ./cmd/main.go


.PHONY: down
down:
	$(DOCKER_COMPOSE_BIN) down -v

.PHONY: up
up: up-docker up-migrate

.PHONY: up-docker
up-docker:
	$(DOCKER_COMPOSE_BIN) up -d

.PHONY: up-migrate
up-migrate:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING='$(POSTGRES_STRING)' goose -dir '$(MIGRATIONS_PATH)' up

.PHONY: generate-sqlc
generate-sqlc:
	go tool sqlc generate
