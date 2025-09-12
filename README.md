# Golang Microservice Architecture

This project demonstrates a microservice architecture implemented in Go, consisting of multiple services that communicate with each other to provide a complete e-commerce solution.

## Services

### Order Service

A RESTful API service responsible for managing orders. It communicates with the Inventory Service to check and reserve product stock.

- **Tech Stack**: Go, Gin, PostgreSQL
- **API Type**: REST
- **Default Port**: 8080

[View Order Service Documentation](./order-service/README.md)

### Inventory Service

A gRPC service responsible for managing product inventory.

- **Tech Stack**: Go, gRPC, PostgreSQL
- **API Type**: gRPC
- **Default Port**: 9090

[View Inventory Service Documentation](./inventory-service/README.md)

## Architecture Overview

The system follows a microservice architecture pattern:

1. **Order Service**: Handles order creation and management
2. **Inventory Service**: Manages product inventory and stock reservation

Services communicate using gRPC for efficient inter-service communication.

### Detailed Architecture

#### System Components

```
┌─────────────────┐      ┌─────────────────┐
│                 │      │                 │
│   Order Service │◄─────┤ Inventory Service│
│    (REST API)   │      │     (gRPC)      │
│                 │      │                 │
└────────┬────────┘      └────────┬────────┘
         │                        │         
         │                        │         
         ▼                        ▼         
┌─────────────────┐      ┌─────────────────┐
│                 │      │                 │
│  Order Database │      │Inventory Database│
│   (PostgreSQL)  │      │   (PostgreSQL)  │
│                 │      │                 │
└─────────────────┘      └─────────────────┘
```

#### Communication Flow

1. **Client Request Flow**:
   - Client sends HTTP request to Order Service
   - Order Service processes the request
   - If inventory check is needed, Order Service calls Inventory Service via gRPC
   - Order Service responds to client with HTTP response

2. **Inter-Service Communication**:
   - Order Service acts as a gRPC client
   - Inventory Service acts as a gRPC server
   - Protocol Buffers are used for data serialization
   - Service discovery is handled via environment variables

#### Design Patterns

- **Database per Service**: Each microservice has its own database
- **API Gateway**: Order Service acts as an entry point for clients
- **Service-to-Service Communication**: gRPC is used for efficient communication between services
- **Event-Driven Architecture**: Services communicate asynchronously when appropriate

## Database

Each service has its own PostgreSQL database to ensure service independence and data isolation.

## Running the Project

### Prerequisites

- Go 1.16+
- Docker and Docker Compose
- PostgreSQL

### Using Docker Compose

```bash
# Start all services and databases
docker-compose up

# Start specific service
docker-compose up order-service
```

#### Auto Seeder

The docker-compose configuration includes auto seeding for the inventory service. When you start the services using docker-compose, the inventory database will be automatically populated with sample data. This feature is enabled by setting the `RUN_SEEDER=true` environment variable in the docker-compose.yml file.

To disable auto seeding, edit the docker-compose.yml file and set `RUN_SEEDER=false` or remove the environment variable.

### Running Locally

```bash
# Start Order Service
cd order-service
go run cmd/main.go

# Start Inventory Service
cd inventory-service
go run cmd/main.go
```

## Development

```bash
# Run tests for all services
make test

# Build all services
make build

# Seed inventory database with sample data
make seed-inventory
```

## Project Structure

```
├── order-service/       # Order management service
├── inventory-service/   # Inventory management service
├── proto/               # Protocol Buffers definitions
├── docker-compose.yml   # Docker Compose configuration
└── Makefile             # Build and development commands
```