.PHONY: up down watch restart clean

up:
	docker compose up -d

down:
	docker compose down

restart:
	docker compose up -d --build

watch:
	docker compose watch

clean:
	docker compose down -v
	rm -rf qdrant_data
