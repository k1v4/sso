all: m start test

m:
	go run ./cmd/migrator --storage-path=./storage/sso.db --migrations-path=./tests/migrations --migrations-table=migrations_tests

start:
	@echo "start app"
	go run ./cmd/sso/main.go --config=./config/local.yaml

test:
	go test -v ./tests