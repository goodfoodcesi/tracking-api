build:
	docker build .
run:
	docker compose up -d --build
stop:
	docker compose down
logs:
	docker logs -f tracking-api-trackingapi-1