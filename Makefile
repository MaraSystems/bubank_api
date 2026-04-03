postgres:
	docker run --name bubank_pg -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine3.21

createdb:
	docker exec -it bubank_pg createdb --username=root --owner=root bubank

dropdb:
	docker exec -it bubank_pg dropdb bubank

migrate_create:
	migrate create -ext sql -dir db/migrations -seq ${name}

migrate_up:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/bubank?sslmode=disable" -verbose up

migrate_down:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/bubank?sslmode=disable" -verbose down

migrate_up_next:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/bubank?sslmode=disable" -verbose up 1

migrate_down_next:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/bubank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

server:
	CompileDaemon -command="./bubank_api"

mock:
	mockgen -destination db/mock/store.go -package mockdb github.com/MaraSystems/bubank_api/db/sqlc Store

test:
	go test -v -cover ./...

document:
	swag fmt & swag init .

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
    proto/*.proto

evans:
	evans --port 9090 -r repl --host localhost --package pb

.PHONY: postgres proto