images:
	sudo docker pull postgres
container:
	sudo docker run --name banking -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root123 -d postgres
createdb:
	sudo docker exec -it banking createdb --username=root --owner=root bank
dropdb:
	sudo docker exec -it banking dropdb bank
migrateup:
	migrate -path db/migration -database "postgresql://root:root123@localhost:5432/bank?sslmode=disable" -verbose up
migrateup1:
	migrate -path db/migration -database "postgresql://root:root123@localhost:5432/bank?sslmode=disable" -verbose up 1
migratedown:
	migrate -path db/migration -database "postgresql://root:root123@localhost:5432/bank?sslmode=disable" -verbose down
migratedown1:
	migrate -path db/migration -database "postgresql://root:root123@localhost:5432/bank?sslmode=disable" -verbose down 1
migrationcreate:
	migrate create -ext sql -dir db/migration -seq add_users
sqlcinit:
	sqlc init
sqlc:
	sqlc generate	
test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/amallick86/Go_bank/db/sqlc Store

.PHONY: images createcontainer createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlcinit sqlc test server mockgen migrationcreate