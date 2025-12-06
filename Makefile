BACKEND_PATH=./backend
API_PATH=./api
GEN_PATH=internal/pb
MIGRATIONS_PATH=migrations
POSTGRES_STRING=postgresql://user:password@localhost:5432/chattery?sslmode=disable



.PHONY: run
run:
	go tool air -c .air.toml

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

.PHONY: down-docker
down-docker:
	docker-compose down -v

.PHONY: up-docker
up-docker:
	docker-compose up -d

.PHONY: generate
generate: generate-user-service

.PHONY: generate-user-service
generate-user-service:
	protoc -I '$(API_PATH)' \
    	--go_out '$(BACKEND_PATH)/user_service/$(GEN_PATH)' \
    	--go_opt paths=source_relative \
    	--go-grpc_out '$(BACKEND_PATH)/user_service/$(GEN_PATH)' \
    	--go-grpc_opt paths=source_relative \
    	--grpc-gateway_out '$(BACKEND_PATH)/user_service/$(GEN_PATH)' \
    	--grpc-gateway_opt paths=source_relative \
    	--openapiv2_out '$(BACKEND_PATH)/user_service/$(GEN_PATH)' \
        '$(API_PATH)/user_service/user_service.proto'

.PHONY: get-grpc-deps
get-grpc-deps:
	rm -rf '$(API_PATH)/google/api'
	wget -qP '$(API_PATH)/google/api' https://raw.githubusercontent.com/googleapis/googleapis/refs/heads/master/google/api/annotations.proto
	wget -qP '$(API_PATH)/google/api' https://raw.githubusercontent.com/googleapis/googleapis/refs/heads/master/google/api/field_behavior.proto
	wget -qP '$(API_PATH)/google/api' https://raw.githubusercontent.com/googleapis/googleapis/refs/heads/master/google/api/http.proto
	wget -qP '$(API_PATH)/google/api' https://raw.githubusercontent.com/googleapis/googleapis/refs/heads/master/google/api/httpbody.proto

.PHONY: migrate
migrate: migrate-user-service

.PHONY: migrate-user-service
migrate-user-service:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING='$(POSTGRES_STRING)' goose -dir '$(BACKEND_PATH)/user_service/$(MIGRATIONS_PATH)' up

.PHONY: generate-sqlc
generate-sqlc:
	cd '$(BACKEND_PATH)/user_service' && sqlc generate
