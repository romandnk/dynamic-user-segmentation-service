run:
	docker compose -f ./deployments/docker-compose.yaml up -d --build

stop:
	docker compose -f ./deployments/docker-compose.yaml down

clear:
	docker volume rm dynamic-user-segmentation_pgdata