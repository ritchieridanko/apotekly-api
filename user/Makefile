# === Variables ===
BINARY_DIR := bin
SERVER_BIN := $(BINARY_DIR)/server
MIGRATE_BIN := $(BINARY_DIR)/migrate

# === Default Target ===
help:
	@echo "Available commands:"
	@echo "  make dev-up				Create the database, apply all migrations, and run the server"
	@echo "  make dev-reset			Drop and create new the database, apply all migrations, and run the server"
	@echo "  make build-and-run			Build and run the migrate and server binaries"
	@echo "  make run-server			Run the server locally"
	@echo "  make build-server			Build the server binary into $(SERVER_BIN)"
	@echo "  make build-and-run-server		Build and run the server binary"
	@echo "  make setup-db				Create the database"
	@echo "  make drop-db				Drop the database"
	@echo "  make migrate-up			Apply all migrations"
	@echo "  make migrate-down			Rollback N migrations"
	@echo "  make migrate-down-all			Rollback all migrations"
	@echo "  make build-migrate			Build the migrate binary into $(MIGRATE_BIN)"
	@echo "  make build-and-run-migrate		Build and run the migrate binary"

# === Server Commands ===
run-server:
	go run cmd/server/main.go

build-server:
	go build -o $(SERVER_BIN) cmd/server/main.go

build-and-run-server:
	make build-server
	./$(SERVER_BIN)

# === DB Commands ===
setup-db:
	psql -U postgres -h localhost -tc "SELECT 1 FROM pg_database WHERE datname = 'apotekly_user_db'" | grep -q 1 || \
	psql -U postgres -h localhost -c "CREATE DATABASE apotekly_user_db"

drop-db:
	psql -U postgres -h localhost -c "DROP DATABASE IF EXISTS apotekly_user_db"

# === Migration Commands ===
migrate-up:
	go run cmd/migrate/main.go -up

migrate-down:
	go run cmd/migrate/main.go -down $(steps)

migrate-down-all:
	go run cmd/migrate/main.go -down 0

build-migrate:
	go build -o $(MIGRATE_BIN) cmd/migrate/main.go

build-and-run-migrate:
	make build-migrate
	./$(MIGRATE_BIN) -up

# === Dev Commands ===
dev-up:
	make setup-db
	make migrate-up
	make run-server

dev-reset:
	make drop-db
	make setup-db
	make migrate-up
	make run-server

build-and-run:
	make drop-db
	make setup-db
	make build-and-run-migrate
	make build-and-run-server