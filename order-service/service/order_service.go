package service

import (
	"context"
	"fmt"
	"log"

	"github.com/fardannozami/golang-microservice/order-service/repository"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	// OrderStatusPending represents a pending order
	OrderStatusPending OrderStatus = "pending"
	// OrderStatusConfirmed represents a confirmed order
	OrderStatusConfirmed OrderStatus = "confirmed"
	// OrderStatusRejected represents a rejected order
	OrderStatusRejected OrderStatus = "rejected"
)

// CreateOrderRequest represents a request to create an order
type CreateOrderRequest struct {
	UserID string
	Items  []OrderItemRequest
}

// OrderItemRequest represents a request to create an order item
type OrderItemRequest struct {
	ProductID string
	Quantity  int
	Price     float64
}

// OrderService defines the interface for order service operations
type OrderService interface {
	CreateOrder(ctx context.Context, req *CreateOrderRequest) (*repository.Order, error)
	GetOrder(ctx context.Context, id string) (*repository.Order, error)
	ListOrders(ctx context.Context) ([]*repository.Order, error)
}

// orderService implements OrderService interface
type orderService struct {
	orderRepo       repository.OrderRepository
	inventoryClient InventoryClient
}

// NewOrderService creates a new order service
func NewOrderService(orderRepo repository.OrderRepository, inventoryClient InventoryClient) OrderService {
	return &orderService{
		orderRepo:       orderRepo,
		inventoryClient: inventoryClient,
	}
}

// CreateOrder creates a new order
func (s *orderService) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*repository.Order, error) {
	// Validate request
	if err := validateCreateOrderRequest(req); err != nil {
		return nil, err
	}

	// Do not pre-check inventory to avoid TOCTOU; rely on atomic reservation

	// Create order
	order := &repository.Order{
		UserID: req.UserID,
		Status: string(OrderStatusPending),
		Items:  make([]repository.OrderItem, len(req.Items)),
	}

	// Convert order items
	for i, item := range req.Items {
		order.Items[i] = repository.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	// Create order in database
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Reserve inventory for all items
	var reservationErrors []error
	for _, item := range order.Items {
		log.Printf("[order-service] Reserving stock product_id=%s qty=%d order_id=%s", item.ProductID, item.Quantity, order.ID)
		err := s.inventoryClient.ReserveStock(ctx, item.ProductID, item.Quantity, order.ID)
		if err != nil {
			reservationErrors = append(reservationErrors, err)
		}
	}

	// If any reservation failed, release all reservations and reject order
	if len(reservationErrors) > 0 {
		// Release all successful reservations
		for _, item := range order.Items {
			_ = s.inventoryClient.ReleaseStock(ctx, item.ProductID, item.Quantity, order.ID)
		}

		// Update order status to rejected
		order.Status = string(OrderStatusRejected)
		_ = s.orderRepo.Update(ctx, order)

		return nil, fmt.Errorf("failed to reserve inventory: %v", reservationErrors[0])
	}

	// Update order status to confirmed
	order.Status = string(OrderStatusConfirmed)
	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	return order, nil
}

// GetOrder gets an order by ID
func (s *orderService) GetOrder(ctx context.Context, id string) (*repository.Order, error) {
	return s.orderRepo.GetByID(ctx, id)
}

// ListOrders lists all orders
func (s *orderService) ListOrders(ctx context.Context) ([]*repository.Order, error) {
	return s.orderRepo.List(ctx)
}

// validateCreateOrderRequest validates a create order request
func validateCreateOrderRequest(req *CreateOrderRequest) error {
	// Check if user ID is provided
	if req.UserID == "" {
		return fmt.Errorf("user ID is required")
	}

	// Check if items are provided
	if len(req.Items) == 0 {
		return fmt.Errorf("at least one item is required")
	}

	// Validate each item
	for i, item := range req.Items {
		// Check if product ID is provided
		if item.ProductID == "" {
			return fmt.Errorf("product ID is required for item %d", i)
		}

		// Check if quantity is valid
		if item.Quantity <= 0 {
			return fmt.Errorf("quantity must be positive for item %d", i)
		}

		// Check if price is valid
		if item.Price <= 0 {
			return fmt.Errorf("price must be positive for item %d", i)
		}
	}

	return nil
}
