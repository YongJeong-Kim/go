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
	migrate -path migration -database "postgresql://root:1234@localhost:5432/mytest_db?sslmode=disable" -verbose up

migratedown:
	migrate -path migration -database "postgresql://root:1234@localhost:5432/mytest_db?sslmode=disable" -verbose down

.PHONY: postgres createdb dropdb migrateup migratedown