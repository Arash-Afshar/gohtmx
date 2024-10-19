setup:
	npm install
	npm install -D tailwindcss
	go install github.com/air-verse/air@latest
	go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

dev:
	air -c .air.toml

dev-css:
	npx tailwindcss -i assets/styles.css -o static/styles.css --postcss --watch

# Example: make create_migration name=init
create_migration:
	migrate create -ext sql -dir pkg/db/migrations -seq $(name)

migrate_up:
	migrate -path pkg/db/migrations -database "sqlite3://db.sqlite" -verbose up

# Example: make migrate_down count=1 to rollback the last migration
migrate_down:
	migrate -path pkg/db/migrations -database "sqlite3://db.sqlite" -verbose down $(count)

