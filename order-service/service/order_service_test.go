package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/fardannozami/golang-microservice/order-service/repository"
	"github.com/fardannozami/golang-microservice/order-service/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOrderRepository is a mock implementation of OrderRepository
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(ctx context.Context, order *repository.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) GetByID(ctx context.Context, id string) (*repository.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.Order), args.Error(1)
}

func (m *MockOrderRepository) List(ctx context.Context) ([]*repository.Order, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*repository.Order), args.Error(1)
}

func (m *MockOrderRepository) Update(ctx context.Context, order *repository.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

// MockInventoryClient is a mock implementation of InventoryClient
type MockInventoryClient struct {
	mock.Mock
}

func (m *MockInventoryClient) CheckStock(ctx context.Context, productID string, quantity int) (bool, error) {
	args := m.Called(ctx, productID, quantity)
	return args.Bool(0), args.Error(1)
}

func (m *MockInventoryClient) ReserveStock(ctx context.Context, productID string, quantity int, orderID string) error {
	args := m.Called(ctx, productID, quantity, orderID)
	return args.Error(0)
}

func (m *MockInventoryClient) ReleaseStock(ctx context.Context, productID string, quantity int, orderID string) error {
	args := m.Called(ctx, productID, quantity, orderID)
	return args.Error(0)
}

func (m *MockInventoryClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestCreateOrder_Success(t *testing.T) {
	// Create mocks
	orderRepo := new(MockOrderRepository)
	inventoryClient := new(MockInventoryClient)

	// Create service
	orderService := service.NewOrderService(orderRepo, inventoryClient)

	// Create request
	req := &service.CreateOrderRequest{
		UserID: "user123",
		Items: []service.OrderItemRequest{
			{
				ProductID: "product123",
				Quantity:  2,
				Price:     10.0,
			},
		},
	}

	// Set up expectations
	inventoryClient.On("CheckStock", mock.Anything, "product123", 2).Return(true, nil)
	orderRepo.On("Create", mock.Anything, mock.AnythingOfType("*repository.Order")).Return(nil)
	inventoryClient.On("ReserveStock", mock.Anything, "product123", 2, mock.AnythingOfType("string")).Return(nil)
	orderRepo.On("Update", mock.Anything, mock.AnythingOfType("*repository.Order")).Return(nil)

	// Call service
	order, err := orderService.CreateOrder(context.Background(), req)

	// Assert expectations
	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, "user123", order.UserID)
	assert.Equal(t, string(service.OrderStatusConfirmed), order.Status)
	assert.Len(t, order.Items, 1)
	assert.Equal(t, "product123", order.Items[0].ProductID)
	assert.Equal(t, 2, order.Items[0].Quantity)
	assert.Equal(t, 10.0, order.Items[0].Price)

	// Verify mocks
	orderRepo.AssertExpectations(t)
	inventoryClient.AssertExpectations(t)
}

func TestCreateOrder_InventoryUnavailable(t *testing.T) {
	// Create mocks
	orderRepo := new(MockOrderRepository)
	inventoryClient := new(MockInventoryClient)

	// Create service
	orderService := service.NewOrderService(orderRepo, inventoryClient)

	// Create request
	req := &service.CreateOrderRequest{
		UserID: "user123",
		Items: []service.OrderItemRequest{
			{
				ProductID: "product123",
				Quantity:  2,
				Price:     10.0,
			},
		},
	}

	// Set up expectations
	inventoryClient.On("CheckStock", mock.Anything, "product123", 2).Return(false, nil)

	// Call service
	order, err := orderService.CreateOrder(context.Background(), req)

	// Assert expectations
	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Contains(t, err.Error(), "not available")

	// Verify mocks
	orderRepo.AssertExpectations(t)
	inventoryClient.AssertExpectations(t)
}

func TestCreateOrder_ReservationFailed(t *testing.T) {
	// Create mocks
	orderRepo := new(MockOrderRepository)
	inventoryClient := new(MockInventoryClient)

	// Create service
	orderService := service.NewOrderService(orderRepo, inventoryClient)

	// Create request
	req := &service.CreateOrderRequest{
		UserID: "user123",
		Items: []service.OrderItemRequest{
			{
				ProductID: "product123",
				Quantity:  2,
				Price:     10.0,
			},
		},
	}

	// Set up expectations
	inventoryClient.On("CheckStock", mock.Anything, "product123", 2).Return(true, nil)
	orderRepo.On("Create", mock.Anything, mock.AnythingOfType("*repository.Order")).Return(nil)
	inventoryClient.On("ReserveStock", mock.Anything, "product123", 2, mock.AnythingOfType("string")).Return(errors.New("reservation failed"))
	inventoryClient.On("ReleaseStock", mock.Anything, "product123", 2, mock.AnythingOfType("string")).Return(nil)
	orderRepo.On("Update", mock.Anything, mock.AnythingOfType("*repository.Order")).Return(nil)

	// Call service
	order, err := orderService.CreateOrder(context.Background(), req)

	// Assert expectations
	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Contains(t, err.Error(), "failed to reserve inventory")

	// Verify mocks
	orderRepo.AssertExpectations(t)
	inventoryClient.AssertExpectations(t)
}

func TestGetOrder_Success(t *testing.T) {
	// Create mocks
	orderRepo := new(MockOrderRepository)
	inventoryClient := new(MockInventoryClient)

	// Create service
	orderService := service.NewOrderService(orderRepo, inventoryClient)

	// Create order
	order := &repository.Order{
		ID:     "order123",
		UserID: "user123",
		Status: string(service.OrderStatusConfirmed),
		Items: []repository.OrderItem{
			{
				ID:        "item123",
				OrderID:   "order123",
				ProductID: "product123",
				Quantity:  2,
				Price:     10.0,
			},
		},
	}

	// Set up expectations
	orderRepo.On("GetByID", mock.Anything, "order123").Return(order, nil)

	// Call service
	result, err := orderService.GetOrder(context.Background(), "order123")

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, order, result)

	// Verify mocks
	orderRepo.AssertExpectations(t)
	inventoryClient.AssertExpectations(t)
}

func TestGetOrder_NotFound(t *testing.T) {
	// Create mocks
	orderRepo := new(MockOrderRepository)
	inventoryClient := new(MockInventoryClient)

	// Create service
	orderService := service.NewOrderService(orderRepo, inventoryClient)

	// Set up expectations
	orderRepo.On("GetByID", mock.Anything, "order123").Return(nil, errors.New("order not found"))

	// Call service
	result, err := orderService.GetOrder(context.Background(), "order123")

	// Assert expectations
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "order not found")

	// Verify mocks
	orderRepo.AssertExpectations(t)
	inventoryClient.AssertExpectations(t)
}

func TestListOrders_Success(t *testing.T) {
	// Create mocks
	orderRepo := new(MockOrderRepository)
	inventoryClient := new(MockInventoryClient)

	// Create service
	orderService := service.NewOrderService(orderRepo, inventoryClient)

	// Create orders
	orders := []*repository.Order{
		{
			ID:     "order123",
			UserID: "user123",
			Status: string(service.OrderStatusConfirmed),
			Items: []repository.OrderItem{
				{
					ID:        "item123",
					OrderID:   "order123",
					ProductID: "product123",
					Quantity:  2,
					Price:     10.0,
				},
			},
		},
		{
			ID:     "order456",
			UserID: "user456",
			Status: string(service.OrderStatusPending),
			Items: []repository.OrderItem{
				{
					ID:        "item456",
					OrderID:   "order456",
					ProductID: "product456",
					Quantity:  1,
					Price:     20.0,
				},
			},
		},
	}

	// Set up expectations
	orderRepo.On("List", mock.Anything).Return(orders, nil)

	// Call service
	result, err := orderService.ListOrders(context.Background())

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, orders, result)
	assert.Len(t, result, 2)

	// Verify mocks
	orderRepo.AssertExpectations(t)
	inventoryClient.AssertExpectations(t)
}

func TestListOrders_Error(t *testing.T) {
	// Create mocks
	orderRepo := new(MockOrderRepository)
	inventoryClient := new(MockInventoryClient)

	// Create service
	orderService := service.NewOrderService(orderRepo, inventoryClient)

	// Set up expectations
	orderRepo.On("List", mock.Anything).Return(nil, errors.New("database error"))

	// Call service
	result, err := orderService.ListOrders(context.Background())

	// Assert expectations
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")

	// Verify mocks
	orderRepo.AssertExpectations(t)
	inventoryClient.AssertExpectations(t)
}
