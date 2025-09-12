gen:
	@protoc \
		--go_out=. --go-grpc_out=. proto/*.proto

seed-inventory:
	@cd inventory-service && go run cmd/seed/main.go

swagger-order:
	@cd order-service && go install github.com/swaggo/swag/cmd/swag@latest && swag init -g cmd/main.go -o ./docs
