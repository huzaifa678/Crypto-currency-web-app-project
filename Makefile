createdb:
	docker exec -it postgres createdb --username=root --owner=root crypto_db

dropdb:
	docker exec -it postgres dropdb crypto_db

postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

migrateup:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/crypto_db?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/crypto_db?sslmode=disable" -verbose down

sqlc:
	sqlc generate

server:
	go run main.go

test:
	go test -v -cover -short ./...
.PHONY: createdb dropdb postgres migrateup migratedown sqlc