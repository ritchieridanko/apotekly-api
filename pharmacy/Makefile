# === Variables ===
BINARY_DIR := bin
MIGRATE_BIN := $(BINARY_DIR)/migrate

# === Default Target ===
help:
	@echo "Available commands:"
	@echo "  make setup-db				Create the database"
	@echo "  make drop-db				Drop the database"
	@echo "  make migrate-up			Apply all migrations"
	@echo "  make migrate-down			Rollback N migrations"
	@echo "  make migrate-down-all			Rollback all migrations"
	@echo "  make build-migrate			Build the migrate binary into $(MIGRATE_BIN)"
	@echo "  make build-and-run-migrate		Build and run the migrate binary"

# === DB Commands ===
setup-db:
	psql -U postgres -h localhost -tc "SELECT 1 FROM pg_database WHERE datname = 'apotekly_pharmacy_db'" | grep -q 1 || \
	psql -U postgres -h localhost -c "CREATE DATABASE apotekly_pharmacy_db"

drop-db:
	psql -U postgres -h localhost -c "DROP DATABASE IF EXISTS apotekly_pharmacy_db"

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