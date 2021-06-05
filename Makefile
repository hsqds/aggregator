lint:
	golangci-lint run ./...

test:
	go test -race -cover ./internal/...

rss-update:
	go run ./cmd/update

wire:
	wire gen ./internal/...

build-update:
	go build -ldflags '-s -w' -o ./build/rss-reader ./cmd/update

migrate:
	tern migrate --migrations ./migrations

install-tern:
	go get github.com/jackc/tern

install-tools: install-tern
	go get github.com/google/wire/cmd/wire

psql:
	docker-compose exec postgres psql -U postgres