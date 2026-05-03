MIGRATIONS_PATH=migrations
E2E_MIGRATIONS_PATH=e2e/migrations
POSTGRES_STRING=postgresql://user:password@localhost:5432/chattery?sslmode=disable
DOCKER_COMPOSE_BIN=docker-compose
LOCAL_COMPOSE_FILE=docker-compose.local.yml

.PHONY: run
run:
	go tool air -c .air.toml

.PHONY: run-web
run-web:
	cd web && npm run dev

.PHONY: build
build:
	go build -o ./bin/chattery ./cmd/main.go

.PHONY: build-web
build-web:
	cd web && npm run build

.PHONY: down
down:
	$(DOCKER_COMPOSE_BIN) -f $(LOCAL_COMPOSE_FILE) down -v

.PHONY: up
up: up-docker up-migrate up-e2e

.PHONY: up-docker
up-docker:
	$(DOCKER_COMPOSE_BIN) -f $(LOCAL_COMPOSE_FILE) up -d

.PHONY: up-migrate
up-migrate:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING='$(POSTGRES_STRING)' goose -dir '$(MIGRATIONS_PATH)' up

.PHONY: up-e2e
up-e2e:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING='$(POSTGRES_STRING)' goose -dir '$(E2E_MIGRATIONS_PATH)' up

.PHONY: generate-sqlc
generate-sqlc:
	go tool sqlc generate

.PHONY: lint
lint: install-golangci
	./bin/golangci-lint run ./...

.PHONY: test
test:
	go test ./...

.PHONY: install-golangci
install-golangci:
	GOBIN=$(PWD)/bin go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.4
