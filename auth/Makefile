# === Variables ===
BINARY_DIR := bin
SERVER_BIN := $(BINARY_DIR)/server
MIGRATE_BIN := $(BINARY_DIR)/migrate

# === Default Target ===
help:
	@echo "Available commands:"
	@echo "  make run-server			Run the server locally"
	@echo "  make build-server			Build the server binary into $(SERVER_BIN)"
	@echo "  make migrate-up			Apply all migrations"
	@echo "  make migrate-down-all			Rollback all migrations"
	@echo "  make migrate-down			Rollback N migrations"
	@echo "  make build-migrate			Build the migrate binary into $(MIGRATE_BIN)"

# === Server Commands ===
run-server:
	go run cmd/server/main.go

build-server:
	go build -o $(SERVER_BIN) cmd/server/main.go

# === Migration Commands ===
migrate-up:
	go run cmd/migrate/main.go -up

migrate-down:
	go run cmd/migrate/main.go -down $(steps)

migrate-down-all:
	go run cmd/migrate/main.go -down 0

build-migrate:
	go build -o $(MIGRATE_BIN) cmd/migrate/main.go