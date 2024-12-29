all: m start

m:
	go run ./cmd/migrator --storage-path=./storage/sso.db --migrations-path=./migrations

start:
	echo "start app"
	go run ./cmd/sso/main.go --config=./config/local.yaml
