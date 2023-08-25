FROM migrate/migrate:latest

COPY ./migrations /migrations

ENTRYPOINT ["/bin/sh", "-c", "migrate -path ./migrations -database \"$DATABASE_URL\" up"]
