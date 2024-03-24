postgres:
	pg_ctl start

createdb:
	psql -U postgres -c "create database simple_bank"

dropdb:
	psql -U postgres -d <database_name> -c "drop database simple_bank"

migrateup:
	migrate -path db/migration -database "postgresql://postgres:@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test
