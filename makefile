# Postgres variables
PG_CONTAINER = postgres-alpine
PG_IMAGE = postgres:alpine
PG_USER = root
PG_PASSWORD = password
PG_PORT = 5432
PG_DB = simple_bank

postgres:
	@if [ "$$(docker ps -q -f name=$(PG_CONTAINER))" ]; then \
		echo "✅ Postgres container is already running."; \
	elif [ "$$(docker ps -aq -f status=exited -f name=$(PG_CONTAINER))" ]; then \
		echo "▶ Starting existing Postgres container..."; \
		docker start $(PG_CONTAINER); \
	else \
		echo "🚀 Creating and starting Postgres container..."; \
		docker run --name $(PG_CONTAINER) \
			-p $(PG_PORT):5432 \
			-e POSTGRES_USER=$(PG_USER) \
			-e POSTGRES_PASSWORD=$(PG_PASSWORD) \
			-d $(PG_IMAGE); \
	fi

create_db: postgres
	docker exec -it $(PG_CONTAINER) createdb --username=$(PG_USER) --owner=$(PG_USER) $(PG_DB)

drop_db:
	docker exec -it $(PG_CONTAINER) dropdb --username=$(PG_USER) --if-exists $(PG_DB)

migrate_up:
	migrate -path database/migration -database "postgresql://$(PG_USER):$(PG_PASSWORD)@localhost:$(PG_PORT)/$(PG_DB)?sslmode=disable" -verbose up

migrate_down:
	migrate -path database/migration -database "postgresql://$(PG_USER):$(PG_PASSWORD)@localhost:$(PG_PORT)/$(PG_DB)?sslmode=disable" -verbose down

migrate_fresh:
	migrate -path database/migration -database "postgresql://$(PG_USER):$(PG_PASSWORD)@localhost:$(PG_PORT)/$(PG_DB)?sslmode=disable" -verbose down && \
	migrate -path database/migration -database "postgresql://$(PG_USER):$(PG_PASSWORD)@localhost:$(PG_PORT)/$(PG_DB)?sslmode=disable" -verbose up

sql_gen:
	sqlc generate

test:
	go test -v -cover ./...

pg-stop:
	docker stop $(PG_CONTAINER)

pg-down:
	docker rm -f $(PG_CONTAINER)

.PHONY: postgres create_db drop_db migrate_up migrate_down migrate_fresh sql_gen test pg-stop pg-down
