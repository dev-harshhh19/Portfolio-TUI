# Harshad Nikam - TUI Portfolio Makefile
.PHONY: all build run ssh test docker clean

all: build

build:
	@echo "==> Building binaries..."
	@mkdir -p bin
	go build -ldflags="-s -w" -o bin/harshad-tui ./cmd/cli
	go build -ldflags="-s -w" -o bin/harshad-ssh ./cmd/ssh
	@echo "==> Build complete: bin/harshad-tui and bin/harshad-ssh"

run:
	go run ./cmd/cli

ssh:
	SSH_PORT=2222 go run ./cmd/ssh

test:
	go test -v ./internal/...

docker-up:
	docker compose -f deployments/docker-compose.yml up -d --build

docker-down:
	docker compose -f deployments/docker-compose.yml down

clean:
	rm -rf bin/
