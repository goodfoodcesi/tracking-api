build:
	docker build .
run:
	docker compose -f docker-compose-rabbitmq.yaml -p rabbitmq up  -d
	docker compose -f docker-compose-redis.yaml -p redis up -d
	docker compose up -d --build
	docker compose -f docker-compose-traefik.yaml -p traefik up -d
stop:
	docker compose down
logs:
	docker logs -f tracking-api-trackingapi-1