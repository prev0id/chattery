API_PATH=./api
MIGRATIONS_PATH=migrations
POSTGRES_STRING=postgresql://user:password@localhost:5432/chattery?sslmode=disable
DOCKER_COMPOSE_BIN=docker compose

.PHONY: run
run:
	go tool air -c .air.toml

.PHONY: build
build:
	npx @tailwindcss/cli -i ./website/app.css -o ./website/src/css/main.css --minify
	go tool templ generate
	go build -o ./bin/chattery ./cmd/main.go

.PHONY: start-tailwind
start-tailwind:
	npx @tailwindcss/cli \
	    -i ./website/app.css \
		-o ./website/src/css/main.css \
		--watch --minify

.PHONY: generate-proto
generate-proto:
	protoc \
		--proto_path='api' \
		--go_out='internal/pb' \
		websocket/websocket.proto

.PHONY: down
down:
	$(DOCKER_COMPOSE_BIN) down -v

.PHONY: up
up: up-docker up-migrate

.PHONY: up-docker
up-docker:
	$(DOCKER_COMPOSE_BIN) up -d

.PHONY: generate
generate:
	protoc -I '$(API_PATH)' \
    	--go_out '$(GEN_PATH)' \
    	--go_opt paths=source_relative \
    	--go-grpc_out '$(GEN_PATH)' \
    	--go-grpc_opt paths=source_relative \
    	--grpc-gateway_out '$(GEN_PATH)' \
    	--grpc-gateway_opt paths=source_relative \
    	--openapiv2_out '$(GEN_PATH)' \
        '$(API_PATH)/user_service/user_service.proto'

.PHONY: get-grpc-deps
get-grpc-deps:
	rm -rf '$(API_PATH)/google/api'
	wget -qP '$(API_PATH)/google/api' https://raw.githubusercontent.com/googleapis/googleapis/refs/heads/master/google/api/annotations.proto
	wget -qP '$(API_PATH)/google/api' https://raw.githubusercontent.com/googleapis/googleapis/refs/heads/master/google/api/field_behavior.proto
	wget -qP '$(API_PATH)/google/api' https://raw.githubusercontent.com/googleapis/googleapis/refs/heads/master/google/api/http.proto
	wget -qP '$(API_PATH)/google/api' https://raw.githubusercontent.com/googleapis/googleapis/refs/heads/master/google/api/httpbody.proto

.PHONY: up-migrate
up-migrate:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING='$(POSTGRES_STRING)' goose -dir '$(MIGRATIONS_PATH)' up

.PHONY: generate-sqlc
generate-sqlc:
	go tool sqlc generate
