server:
	cd cmd/api && go run .

build:
	docker-compose up --build -d

down:
	docker-compose down

.PHONY: server build down