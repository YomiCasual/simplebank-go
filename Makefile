migrateup:
	migrate -path db/migration -database "postgresql://user:password@localhost:5432/simplebank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://user:password@localhost:5432/simplebank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://user:password@localhost:5432/simplebank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://user:password@localhost:5432/simplebank?sslmode=disable" -verbose down 1

migrateuptest:
	migrate -path db/migration -database "postgresql://user:password@localhost:5432/simplebanktest?sslmode=disable" -verbose up

migratedowntest:
	migrate -path db/migration -database "postgresql://user:password@localhost:5432/simplebanktest?sslmode=disable" -verbose down

sqlc: 
	sqlc generate

runtest:
	cd db/sqlc/test && go test .

test:
	go test -v -cover ./...

MIGRATION_NAME ?= $(shell bash -c 'read -p "MigrationName: " name; echo $$name')
DB_NAME ?= $(shell bash -c 'read -p "DbName: " name; echo $$name')

gmigration:
	migrate create -ext sql -dir db/migration -seq $(MIGRATION_NAME)

dockerup:
	docker compose up $(DB_NAME)



.PHONY: migrateup migratedown sqlc runtest test server migrateuptest migratedowntest gmigration migratedown1 migrateup1 dockerup