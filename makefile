include .env

MIGRATE_PATH = cmd/migrate/migrations

DATABASE_NAME=gobali_db

.PHONY: new-migration
migrate:
	@migrate create -ext sql -dir $(MIGRATE_PATH) -seq $(name)

.PHONY:migration-up
migrate-up:
	@migrate -path=./cmd/migrate/migrations -database=mysql://root:$(MYSQL_AUTH)@/$(DATABASE_NAME) up

.PHONY:migration-down
migrate-down:
	@migrate -path=./cmd/migrate/migrations -database=mysql://root:$(MYSQL_AUTH)@/$(DATABASE_NAME) down

.PHONY:migration-back
migrate-back:
	@migrate -path=./cmd/migrate/migrations -database=mysql://root:$(MYSQL_AUTH)@/$(DATABASE_NAME) force $(no)


.PHONY: gen-docs
gen-docs:
	@swag init -g ./api/main.go -d cmd,internal && swag fmt
