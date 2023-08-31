FROM migrate/migrate:latest

COPY ./migrations /migrations

ENTRYPOINT ["/bin/sh", "-c", "migrate -path ./migrations -database \"postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}-user-segmentation_db?sslmode=disable\" up"]
