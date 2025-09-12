# Inventory Service

The Inventory Service is a gRPC service responsible for managing product inventory in the microservice architecture.

## Features

- Check product availability
- Reserve product stock for orders
- Release reserved stock when orders are cancelled
- Manage product inventory levels

## Architecture

The Inventory Service follows a layered architecture:

1. **gRPC Server Layer**: Handles gRPC requests and responses
2. **Service Layer**: Contains business logic
3. **Repository Layer**: Manages database operations

## gRPC API

The service exposes the following gRPC endpoints:

### CheckStock

Checks if a product has sufficient stock available.

```protobuf
rpc CheckStock(CheckStockRequest) returns (CheckStockResponse) {}

message CheckStockRequest {
  string product_id = 1;
  int32 quantity = 2;
}

message CheckStockResponse {
  bool available = 1;
  string message = 2;
}
```

### ReserveStock

Reserves stock for a product when an order is created.

```protobuf
rpc ReserveStock(ReserveStockRequest) returns (ReserveStockResponse) {}

message ReserveStockRequest {
  string product_id = 1;
  int32 quantity = 2;
  string order_id = 3;
}

message ReserveStockResponse {
  bool success = 1;
  string message = 2;
}
```

### ReleaseStock

Releases previously reserved stock when an order is cancelled.

```protobuf
rpc ReleaseStock(ReleaseStockRequest) returns (ReleaseStockResponse) {}

message ReleaseStockRequest {
  string product_id = 1;
  int32 quantity = 2;
  string order_id = 3;
}

message ReleaseStockResponse {
  bool success = 1;
  string message = 2;
}
```

## Database Schema

The service uses two main tables:

### Products

```
+---------------+
| products      |
+---------------+
| id            |
| name          |
| description   |
| price         |
| created_at    |
+---------------+
```

### Inventory

```
+---------------+
| inventory     |
+---------------+
| product_id    |
| quantity      |
| reserved      |
| updated_at    |
+---------------+
```

## Configuration

The service can be configured using environment variables:

- `SERVER_PORT`: gRPC server port (default: 9090)
- `DATABASE_URL`: PostgreSQL connection string

## Running Locally

```bash
# Run the service
go run cmd/main.go

# Run tests
go test ./...

# Seed database with sample data
make seed-inventory
```

## Seeder

The service includes a data seeder to populate the database with sample products and inventory data. The seeder creates:

- 5 sample products (Laptop Gaming, Smartphone, Headphone Bluetooth, Smart Watch, Wireless Earbuds)
- Initial inventory quantities for each product

To run the seeder:

```bash
# From project root
make seed-inventory

# Or directly
cd inventory-service && go run cmd/seed/main.go
```

## Docker

```bash
# Build the Docker image
docker build -t inventory-service .

# Run the container without seeder
docker run -p 9090:9090 \
  -e DATABASE_URL=postgres://postgres:postgres@postgres:5432/inventory_service \
  inventory-service

# Run the container with auto seeder
docker run -p 9090:9090 \
  -e DATABASE_URL=postgres://postgres:postgres@postgres:5432/inventory_service \
  -e RUN_SEEDER=true \
  inventory-service
```

### Auto Seeder in Docker

The Docker image includes an entrypoint script that can automatically run the seeder before starting the service. To enable this feature:

1. Set the environment variable `RUN_SEEDER=true` when running the container
2. The seeder will populate the database with sample data before the service starts
3. This is useful for development and testing environments

**Note:** For production environments, it's recommended to disable auto seeding by setting `RUN_SEEDER=false` or omitting the environment variable.