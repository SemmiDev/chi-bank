migrateup:
	migrate -path db/migration -database "postgresql://sammidev:sammidev@localhost:5432/chi_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://sammidev:sammidev@localhost:5432/chi_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://sammidev:sammidev@localhost:5432/chi_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://sammidev:sammidev@localhost:5432/chi_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: migrateup migratedown migrateup1 migratedown1 sqlc test server
