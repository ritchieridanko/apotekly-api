# === Variables ===
BINARY_DIR := bin
APP_BIN := $(BINARY_DIR)/app
MIGRATE_BIN := $(BINARY_DIR)/migrate
PROTOC := protoc
PROTO_DIR := pkg/events
PROTOC_GEN_GO := $(shell which protoc-gen-go)
PROTOC_GEN_GO_GRPC := $(shell which protoc-gen-go-grpc)

help:
	@echo "Available commands:"
	@echo "	make run-dev				Drop and create the database, apply all migrations, and run the app"
	@echo "	make build-and-run-dev			Build and run migrate and app binaries locally"
	@echo "	make run-app				Run the app locally"
	@echo "	make build-app				Build the app binary"
	@echo "	make build-and-run-app			Build and run the app binary"
	@echo "	make build-protobuf			Build protobuf files"
	@echo "	make drop-protobuf			Drop generated protobuf files"
	@echo "	make setup-database			Create the database (local only)"
	@echo "	make drop-database			Drop the database (local only)"
	@echo "	make migrate-up				Run migrations locally"
	@echo "	make migrate-down			Rollback N migrations locally"
	@echo "	make migrate-down-all			Rollback all migrations locally"
	@echo "	make build-migrate			Build the migrate binary"
	@echo "	make build-and-run-migrate		Build and run migrate locally"
	@echo "	make docker-build			Build Docker containers"
	@echo "	make docker-up				Start Docker containers"
	@echo "	make docker-down			Drop Docker containers"
	@echo "	make docker-migrate-up			Run migrations in Docker container"
	@echo "	make docker-migrate-down		Rollback migrations in Docker container"

# === App Commands ===
run-app:
	make build-protobuf
	go run cmd/app/main.go

build-app:
	make build-protobuf
	go build -o $(APP_BIN) cmd/app/main.go

build-and-run-app:
	make build-app
	./$(APP_BIN)

# === Protobuf Commands ===
build-protobuf:
	make drop-protobuf
	$(PROTOC) \
	--go_out=. --go_opt=paths=source_relative \
	$(PROTO_DIR)/*.proto

drop-protobuf:
	@find $(PROTO_DIR) -name "*.pb.go" -type f -delete

# === Database Commands (local only) ===
setup-database:
	psql -U postgres -h localhost -tc "SELECT 1 FROM pg_database WHERE datname = 'apotekly_auth_db'" | grep -q 1 || \
	psql -U postgres -h localhost -c "CREATE DATABASE apotekly_auth_db"

drop-database:
	psql -U postgres -h localhost -c "DROP DATABASE IF EXISTS apotekly_auth_db"

# === Migration Commands (local only) ===
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

# === Local Dev Commands ===
run-dev:
	make drop-database
	make setup-database
	make migrate-up
	make run-app

build-and-run-dev:
	make drop-database
	make setup-database
	make build-and-run-migrate
	make build-and-run-app

# === Docker Commands ===
docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-migrate-up:
	docker compose run --rm auth-migrate -up

docker-migrate-down:
	docker compose run --rm auth-migrate -down 0
