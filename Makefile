up:
	docker-compose up --build

down:
	docker-compose down

ps:
	docker-compose ps

re: down up

.DEFAULT_GOAL := re
.PHONY: up down ps re