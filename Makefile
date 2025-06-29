server:
	cd cmd/api && go run .

build:
	docker-compose up --build -d

down:
	docker-compose down

proto: 
	rm -f pb/*.go
	protoc \
	--proto_path=proto \
	--go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	proto/*.proto

.PHONY: server build down proto