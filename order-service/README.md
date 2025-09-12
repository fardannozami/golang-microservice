# Order Service

The Order Service is a RESTful API service responsible for managing orders in the microservice architecture. It communicates with the Inventory Service to check and reserve product stock.

## Features

- Create new orders
- List all orders
- Get order details
- Manage order status (pending, confirmed, rejected)
- Communicate with Inventory Service for stock management

## Architecture

The Order Service follows a layered architecture:

1. **HTTP Handler Layer**: Handles HTTP requests and responses using Gin framework
2. **Service Layer**: Contains business logic
3. **Repository Layer**: Manages database operations
4. **Client Layer**: Communicates with other services (Inventory Service)

## REST API

The service exposes the following REST endpoints:

### Create Order

Creates a new order and reserves inventory.

```
POST /api/v1/orders

Request:
{
  "user_id": "user123",
  "items": [
    {
      "product_id": "prod-001",
      "quantity": 2,
      "price": 29.99
    }
  ]
}

Response:
{
  "id": "order123",
  "user_id": "user123",
  "status": "confirmed",
  "items": [
    {
      "id": "item123",
      "product_id": "prod-001",
      "quantity": 2,
      "price": 29.99
    }
  ],
  "created_at": "2023-01-01T12:00:00Z",
  "updated_at": "2023-01-01T12:00:00Z"
}
```

### Get Order

Retrieves an order by ID.

```
GET /api/v1/orders/:id

Response:
{
  "id": "order123",
  "user_id": "user123",
  "status": "confirmed",
  "items": [...],
  "created_at": "2023-01-01T12:00:00Z",
  "updated_at": "2023-01-01T12:00:00Z"
}
```

### List Orders

Retrieves all orders.

```
GET /api/v1/orders

Response:
[
  {
    "id": "order123",
    "user_id": "user123",
    "status": "confirmed",
    "items": [...],
    "created_at": "2023-01-01T12:00:00Z",
    "updated_at": "2023-01-01T12:00:00Z"
  },
  ...
]
```

## Database Schema

The service uses two main tables:

### Orders

```
+---------------+
| orders        |
+---------------+
| id            |
| user_id       |
| status        |
| created_at    |
| updated_at    |
+---------------+
```

### Order Items

```
+---------------+
| order_items   |
+---------------+
| id            |
| order_id      |
| product_id    |
| quantity      |
| price         |
+---------------+
```

## Configuration

The service can be configured using environment variables:

- `SERVER_PORT`: HTTP server port (default: 8080)
- `DATABASE_URL`: PostgreSQL connection string
- `INVENTORY_SERVICE_URL`: URL of the Inventory Service gRPC endpoint

## Running Locally

```bash
# Run the service
go run cmd/main.go

# Run tests
go test ./...
```

## Docker

```bash
# Build the Docker image
docker build -t order-service .

# Run the container
docker run -p 8080:8080 \
  -e DATABASE_URL=postgres://postgres:postgres@postgres:5432/order_service \
  -e INVENTORY_SERVICE_URL=inventory-service:9090 \
  order-service
```