.PHONY: docker-build docker-up seed migrate-up migrate-down

docker-build:
	docker compose up --build

docker-up:
	docker compose up -d --force-recreate

migrate-up:
	docker compose run --rm migrate \
		-path /migrations \
		-database "postgres://postgres:postgres@localhost:5433/booking?sslmode=disable" \
		up

migrate-down:
	docker compose run --rm migrate \
		-path /migrations \
		-database "postgres://postgres:postgres@localhost:5433/booking?sslmode=disable" \
		down --all