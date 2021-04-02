# Linux or powershell
postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1234 -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root mytest_db

dropdb:
	docker exec -it postgres12 dropdb mytest_db

# using git bash
gpostgres:
	winpty docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1234 -d postgres:12-alpine

gcreatedb:
	winpty docker exec -it postgres12 createdb --username=root --owner=root mytest_db

gdropdb:
	winpty docker exec -it postgres12 dropdb mytest_db

migrateup:
	migrate -path db/migration -database "mysql://root:1234@tcp(localhost:13306)/go" -verbose up

migrateup1:
	migrate -path db/migration -database "mysql://root:1234@tcp(localhost:13306)/go" -verbose up 1

migratedown:
	migrate -path db/migration -database "mysql://root:1234@tcp(localhost:13306)/go" -verbose down

migratedown1:
	migrate -path db/migration -database "mysql://root:1234@tcp(localhost:13306)/go" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/yongjeong-kim/go/gogin/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test server mock