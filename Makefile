build:
	go build -o bin/vodeno_app main.go
run:
	docker compose up -d --wait
	bin/vodeno_app
down:
	$(eval PIDAPP := $(shell pidof vodeno_app))
	kill $(PIDAPP)
	docker compose down
