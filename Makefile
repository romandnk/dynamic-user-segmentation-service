make run:
	docker compose -f ./deployments/docker-compose.yaml up -d --build

make stop:
	docker compose -f ./deployments/docker-compose.yaml down