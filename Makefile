.PHONY: run test fmt vet postgres-up mongo-up db-down migrate keygen check
run:
	go run ./cmd/api

test:
	go test ./...
fmt:
	gofmt -w ./cmd ./internal
vet:
	go vet ./...
postgres-up:
	docker compose --profile postgres up -d
mongo-up:
	docker compose --profile mongo up -d
db-down:
	docker compose --profile postgres --profile mongo down
migrate:
	docker compose exec -T postgres psql -U app -d app -v ON_ERROR_STOP=1 < migrations/001_items.sql
keygen:
	go run ./cmd/keygen
check: fmt vet test
