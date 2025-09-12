package service

import (
	"context"
	"fmt"

	"github.com/fardannozami/golang-microservice/inventory-service/repository"
)

// InventoryService defines the interface for inventory service operations
type InventoryService interface {
	CheckStock(ctx context.Context, productID string, quantity int) (bool, error)
	ReserveStock(ctx context.Context, productID string, quantity int, orderID string) error
	ReleaseStock(ctx context.Context, productID string, quantity int, orderID string) error
}

// inventoryService implements InventoryService interface
type inventoryService struct {
	repo repository.InventoryRepository
}

// NewInventoryService creates a new inventory service
func NewInventoryService(repo repository.InventoryRepository) InventoryService {
	return &inventoryService{repo: repo}
}

// CheckStock checks if a product is available in inventory
func (s *inventoryService) CheckStock(ctx context.Context, productID string, quantity int) (bool, error) {
	// Validate input
	if productID == "" {
		return false, fmt.Errorf("product ID is required")
	}
	if quantity <= 0 {
		return false, fmt.Errorf("quantity must be positive")
	}

	// Check stock in repository
	return s.repo.CheckStock(ctx, productID, quantity)
}

// ReserveStock reserves stock for an order
func (s *inventoryService) ReserveStock(ctx context.Context, productID string, quantity int, orderID string) error {
	// Validate input
	if productID == "" {
		return fmt.Errorf("product ID is required")
	}
	if quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}
	if orderID == "" {
		return fmt.Errorf("order ID is required")
	}

	// Directly reserve stock in repository with row-level locking to ensure atomicity
	return s.repo.ReserveStock(ctx, productID, quantity, orderID)
}

// ReleaseStock releases reserved stock
func (s *inventoryService) ReleaseStock(ctx context.Context, productID string, quantity int, orderID string) error {
	// Validate input
	if productID == "" {
		return fmt.Errorf("product ID is required")
	}
	if quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}
	if orderID == "" {
		return fmt.Errorf("order ID is required")
	}

	// Release stock in repository
	return s.repo.ReleaseStock(ctx, productID, quantity, orderID)
}
