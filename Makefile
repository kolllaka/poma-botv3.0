run:
	go run ./cmd/url-shortener/main.go --config=./config/local.yaml

migrate-up:
	go run ./cmd/migrator/main.go --m=up --storage-path=./storage/music.db --migrations-path=./storage/migrations

migrate-down:
	go run ./cmd/migrator/main.go --m=down --storage-path=./storage/music.db --migrations-path=./storage/migrations