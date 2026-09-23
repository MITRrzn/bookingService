.PHONY: docker-build docker-up migrate-up migrate-down

docker-build:
	docker compose up --build -d

docker-up:
	docker compose up -d --force-recreate

migrate-up:
	docker compose run --rm migrate \
		-path /migrations \
		-database "postgres://postgres:postgres@postgres:5432/booking?sslmode=disable" \
		up

migrate-down:
	docker compose run --rm migrate \
		-path /migrations \
		-database "postgres://postgres:postgres@postgres:5432/booking?sslmode=disable" \
		down --all