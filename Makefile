postgres:
	pg_ctl start
docker_postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:14-alpine


createdb:
	psql -U postgres -c "create database simple_bank"

docker_createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank


dropdb:
	psql -U postgres -d simple_bank -c "drop database simple_bank"

migrateup:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test
