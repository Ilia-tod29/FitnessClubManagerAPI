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

migrateup:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/fitness_club_manager?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/fitness_club_manager?sslmode=disable" -verbose down

migrateforce:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/fitness_club_manager?sslmode=disable" -verbose force #add a version after the force

sqlc:
	sqlc generate

.PHONY: postgres stoppostgres rmpostgres createdb dropdb migrateup migratedown sqlc