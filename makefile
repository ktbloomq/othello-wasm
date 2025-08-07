download:
	go mod download
build:
	GOOS=js GOARCH=wasm go build -o web/main.wasm
up:
	docker compose up -d
down:
	docker compose down -d
