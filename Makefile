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
