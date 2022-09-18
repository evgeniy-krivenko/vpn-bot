include .env

.PHONY:
.SILENT:

build:
	go build -o ./.bin/bot cmd/bot/main.go

run: build
	./.bin/bot

check-connections:
	@docker-compose exec pgdb psql -U postgres -c 'SELECT pid, user, state, query FROM pg_stat_activity WHERE state IS NOT NULL;'

create-migration:
	@migrate create -ext sql -dir ./migrations -seq $(name)

migration-up:
	@migrate -path ./migrations -database $(DB_STRING) up

migration-down:
	@migrate -path ./migrations -database $(DB_STRING) down $(n)