postgres13 :
	docker run --name postgres13 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:13-alpine


createdb:
	docker exec -it postgres13 createdb -U root simple_bank

# migrateup will pull all the schema from the sql file into the database 
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

dropdb:
	docker exec -it postgres13 dropdb -U root simple_bank

sqlc:
	sqlc generate

