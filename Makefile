dockerpostgresimages:
	sudo docker pull postgres
createcontainer:
	sudo docker run --name banking -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root123 -d postgres
createdb:
	sudo docker exec -it banking createdb --username=root --owner=root bank
dropdb:
	sudo docker exec -it banking dropdb bank
migrateup:
	migrate -path db/migration -database "postgresql://root:root123@localhost:5432/bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:root123@localhost:5432/bank?sslmode=disable" -verbose down
sqlcinit:
	sqlc init
sqlc:
	sqlc generate	
test:
	go test -v -cover ./...
server:
	go run main.go
mockgen:
	mockgen -package mockdb -destination db/mock/store.go github.com/amallick86/Go_bank/db/sqlc Store

.PHONY: dockerpostgresimages createcontainer createdb dropdb migrateup migratedown sqlcinit sqlc test server mockgen