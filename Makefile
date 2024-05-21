clean:
	docker compose down
buildup:
	docker-compose up --build
up:
	docker compose up
restart:
	make clean && make up
logs:
	docker-compose logs