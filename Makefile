.PHONY: run
run:
	go tool air -c .air.toml

.PHONY: start-tailwind
start-tailwind:
	npx @tailwindcss/cli \
	    -i ./website/app.css \
		-o ./website/src/css/main.css \
		--watch --minify
