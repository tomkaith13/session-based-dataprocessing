clean:
	docker-compose down --rmi all -v
buildup:
	docker-compose up --build
up:
	docker compose up
restart:
	make clean && make up
logs:
	docker-compose logs