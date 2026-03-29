postgres:
	docker run --name graybank_pg -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine3.21

createdb:
	docker exec -it graybank_pg createdb --username=root --owner=root graybank

dropdb:
	docker exec -it graybank_pg dropdb graybank

migrate_create:
	migrate create -ext sql -dir db/migrations -seq ${name}

migrate_up:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/graybank?sslmode=disable" -verbose up

migrate_down:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/graybank?sslmode=disable" -verbose down

migrate_up_next:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/graybank?sslmode=disable" -verbose up 1

migrate_down_next:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/graybank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

serve:
	CompileDaemon -command="./graybank_api"

mock:
	mockgen -destination db/mock/store.go -package mockdb github.com/MaraSystems/graybank_api/db/sqlc Store

test:
	go test -v -cover ./...

document:
	swag fmt & swag init .

.PHONY: postgres