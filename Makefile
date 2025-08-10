createdb:
	docker exec -it postgres createdb --username=root --owner=root crypto_db

dropdb:
	docker exec -it postgres dropdb crypto_db

postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

migrateup:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/crypto_db?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/crypto_db?sslmode=disable" -verbose down

sqlc:
	sqlc generate

server:
	go run main.go

mock:
	mockgen -destination db/mock/store.go github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc Store_interface
	mockgen -package mockwk -destination worker/mock/distributor.go github.com/huzaifa678/Crypto-currency-web-app-project/worker TaskDistributor

test:
	go test -v -cover -short ./...

go-backend-compose:
	docker compose up --build  

proto:
	rm -rf pb/*.go
	rm -f docs/*.swagger.json
	protoc --proto_path=proto-files --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt paths=source_relative \
	--openapiv2_out=docs --openapiv2_opt=allow_merge=true,merge_file_name=Crypto-currency-web-app \
    proto-files/*.proto

evans:
	evans --host localhost --port 9090 -r repl

go-tools:
	go get -tool github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
	go get -tool github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	go get -tool google.golang.org/protobuf/cmd/protoc-gen-go
	go get -tool google.golang.org/grpc/cmd/protoc-gen-go-grpc

redis:
	docker run --name redis -p 6379:6379 -d redis:7.2-alpine

.PHONY: createdb dropdb postgres migrateup migratedown sqlc server mock test go-backend-compose proto evans go-tools redis