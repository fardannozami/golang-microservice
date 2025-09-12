package service

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/fardannozami/golang-microservice/inventory-service/proto"
)

// InventoryClient defines the interface for inventory client operations
type InventoryClient interface {
	CheckStock(ctx context.Context, productID string, quantity int) (bool, error)
	ReserveStock(ctx context.Context, productID string, quantity int, orderID string) error
	ReleaseStock(ctx context.Context, productID string, quantity int, orderID string) error
	Close() error
}

// inventoryClient implements InventoryClient interface
type inventoryClient struct {
	conn   *grpc.ClientConn
	client pb.InventoryServiceClient
}

// NewInventoryClient creates a new inventory client
func NewInventoryClient(inventoryServiceURL string) (InventoryClient, error) {
	// Configure connection parameters with retry
	connParams := grpc.WithConnectParams(grpc.ConnectParams{
		Backoff: backoff.Config{
			BaseDelay:  100 * time.Millisecond,
			Multiplier: 1.6,
			Jitter:     0.2,
			MaxDelay:   3 * time.Second,
		},
		MinConnectTimeout: 5 * time.Second,
	})

	// Create connection
	conn, err := grpc.Dial(
		inventoryServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		connParams,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to inventory service: %w", err)
	}

	// Create client
	client := pb.NewInventoryServiceClient(conn)

	return &inventoryClient{
		conn:   conn,
		client: client,
	}, nil
}

// CheckStock checks if a product is available in inventory
func (c *inventoryClient) CheckStock(ctx context.Context, productID string, quantity int) (bool, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Call inventory service
	resp, err := c.client.CheckStock(ctx, &pb.CheckStockRequest{
		ProductId: productID,
		Quantity:  int32(quantity),
	})
	if err != nil {
		return false, fmt.Errorf("failed to check stock: %w", err)
	}

	return resp.Available, nil
}

// ReserveStock reserves stock for an order
func (c *inventoryClient) ReserveStock(ctx context.Context, productID string, quantity int, orderID string) error {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Call inventory service
	resp, err := c.client.ReserveStock(ctx, &pb.ReserveStockRequest{
		ProductId: productID,
		Quantity:  int32(quantity),
		OrderId:   orderID,
	})
	if err != nil {
		return fmt.Errorf("failed to reserve stock: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("failed to reserve stock: %s", resp.Message)
	}

	return nil
}

// ReleaseStock releases reserved stock
func (c *inventoryClient) ReleaseStock(ctx context.Context, productID string, quantity int, orderID string) error {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Call inventory service
	resp, err := c.client.ReleaseStock(ctx, &pb.ReleaseStockRequest{
		ProductId: productID,
		Quantity:  int32(quantity),
		OrderId:   orderID,
	})
	if err != nil {
		return fmt.Errorf("failed to release stock: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("failed to release stock: %s", resp.Message)
	}

	return nil
}

// Close closes the connection
func (c *inventoryClient) Close() error {
	return c.conn.Close()
}
