GOOSE_DBSTRING ?= "user=root password=secret host=127.0.0.1 port=5678 dbname=simple-bank sslmode=disable"
GOOSE_MIGRATION_DIR ?= db/migration
GOOSE_DRIVER ?= postgres


network:
	docker network create bank-network

postgres:
	docker run --name simple-bank --network bank-network -p 5678:5432 -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=root -d postgres:15-alpine

build-server:
	docker build -t simplebank:latest .

server:
	docker run --name bankserver --network bank-network -p 8002:8002 -e GIN_MODE=release simplebank:latest 

createdb:
	docker exec -it simple-bank createdb --username=root --owner=root simple-bank

dropdb:
	docker exec -it simple-bank dropdb simple-bank

migrate-up:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir=$(GOOSE_MIGRATION_DIR) up


migrate-down:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir=$(GOOSE_MIGRATION_DIR) down to 0

create-migration:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir=$(GOOSE_MIGRATION_DIR) -s create $(name) sql

sqlc:
	sqlc generate

dev:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/phongnd2802/simple-bank/db/sqlc Store

test:
	go test -v -cover ./...
.PHONY: server build-server mock dev test network postgres createdb dropdb migrate-up migrate-down create-migration sqlc
