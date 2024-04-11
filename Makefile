DB_URL=postgresql://postgres:postgres@localhost:5432/simple_bank?sslmode=disable

postgres:
	pg_ctl start

create_network:
	docker network create bank-network

docker_postgres:
	docker run --name postgres --network bank-network -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:14-alpine

createdb:
	psql -U postgres -c "create database simple_bank"

createdb_postgres:
	psql -h postgres -U postgres -c "create database simple_bank"

docker_createdb:
	docker exec -it postgres createdb --username=postgres --owner=postgres simple_bank

dropdb:
	psql -U postgres -d simple_bank -c "drop database simple_bank"

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/vantu-fit/master-go-be/db/sqlc Store
	mockgen -package mockwk -destination worker/mock/distributor.go github.com/vantu-fit/master-go-be/worker TaskDistributor

proto:
	rm -f pb/*.go 
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt paths=source_relative \
	--openapiv2_out=doc/swagger  --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
    proto/*.proto
	statik -src=./doc/swagger -dest=./doc
evans:
	evans --host localhost --port 9090 -r repl

db_docs:
	dbdocs build doc/db.dbml

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine


.PHONY: postgres createdb dropdb migrateup migratedown sqlc test mock proto evans redis
