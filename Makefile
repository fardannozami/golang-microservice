gen:
	@protoc \
		--go_out=. --go-grpc_out=. proto/*.proto

seed-inventory:
	@cd inventory-service && go run cmd/seed/main.go
