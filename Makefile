clean:
	docker compose down
up:
	docker compose up
restart:
	make clean && make up
logs:
	docker-compose logs