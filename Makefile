migrateup:
	migrate -path db/migration -database "postgresql://user:password@localhost:5432/simplebank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://user:password@localhost:5432/simplebank?sslmode=disable" -verbose down

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



.PHONY: migrateup migratedown sqlc runtest test server migrateuptest migratedowntest