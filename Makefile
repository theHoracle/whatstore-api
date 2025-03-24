# ...existing code...

.PHONY: migrate-generate
migrate-generate:
	go run cmd/migrate/main.go $(name)

.PHONY: migrate-up
migrate-up:
	go run main.go

.PHONY: migrate-down
migrate-down:
	migrate -database ${DATABASE_URL} -path db/migrations down