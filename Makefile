include .env
export

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATIONS_DIR=./migrations

.PHONY: migrate-up migrate-down migrate-force migrate-version migrate-create migrate-drop

## Применить все новые миграции
migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

## Откатить последнюю миграцию
migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

## Откатить всё (осторожно)
migrate-drop:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" drop -f

## Создать новую миграцию: make migrate-create name=add_posts_table
migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

## Показать текущую версию миграции
migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

## Принудительно выставить версию (если миграция упала в dirty state): make migrate-force version=1
migrate-force:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $(version)

## Поднять докер-контейнеры
up:
	docker-compose up -d

## Остановить докер-контейнеры
down:
	docker-compose down

## Запустить приложение
run:
	go run cmd/api/main.go