.PHONY: all up migrate generate subscriber database

all: up migrate

up:
	docker-compose up -d

migrate: wait_for_db
	docker-compose exec database sh -c 'psql -U casino < /db/migrations/00001.create_base.sql'
	docker-compose exec database sh -c 'psql -U casino < /db/migrations/00002.create_events.sql'

wait_for_db:
	@echo "Waiting for database to be ready..."
	@until docker-compose exec database pg_isready -U casino > /dev/null 2>&1; do \
		echo "Waiting for database..."; \
		sleep 2; \
	done
	@echo "Database is ready!"

generator:
	docker-compose --profile generator up -d generator

subscriber:
	docker-compose --profile subscriber up -d subscriber

database:
	docker-compose exec database psql -U casino

run: all subscriber generator
	@echo "Starting services..."
	@echo "Metrics endpoint will be available at http://localhost:8080/materialized"
	@echo "Listening for metrics..."