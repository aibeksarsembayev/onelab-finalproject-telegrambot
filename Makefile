build:
	docker-compose up
	docker image prune

stop:
	docker-compose down

delete:
	docker rmi alpine:latest onelab-finalproject-telegrambot_app:latest postgres:latest golang:1.19 