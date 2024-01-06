postgres:
	docker run --name postgresDB -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:latest

stoppostgres:
	docker stop postgresDB

rmpostgres:
	docker rm postgresDB

createdb:
	docker exec -it postgresDB createdb --username=root --owner=root fitness_club_manager

dropdb:
	docker exec -it postgresDB dropdb fitness_club_manager

createmigration:
	migrate create -ext sql -dir db/migration -seq add_sessions #add a name of the migration - e.g. init_schema

migrateup:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/fitness_club_manager?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/fitness_club_manager?sslmode=disable" -verbose down

migrateforce:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/fitness_club_manager?sslmode=disable" -verbose force 1 #add a version after the force - e.g. 1

sqlc:
	sqlc generate

server:
	go run main.go

stripelistener:
	stripe listen --forward-to localhost:8080/webhook

.PHONY: postgres stoppostgres rmpostgres createdb dropdb createmigration migrateup migratedown sqlc server stripelistener