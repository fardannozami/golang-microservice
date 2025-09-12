package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/fardannozami/golang-microservice/inventory-service/repository"
	"github.com/fardannozami/golang-microservice/inventory-service/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockInventoryRepository is a mock implementation of InventoryRepository
type MockInventoryRepository struct {
	mock.Mock
}

func (m *MockInventoryRepository) CheckStock(ctx context.Context, productID string, quantity int) (bool, error) {
	args := m.Called(ctx, productID, quantity)
	return args.Bool(0), args.Error(1)
}

func (m *MockInventoryRepository) ReserveStock(ctx context.Context, productID string, quantity int) error {
	args := m.Called(ctx, productID, quantity)
	return args.Error(0)
}

func (m *MockInventoryRepository) ReleaseStock(ctx context.Context, productID string, quantity int) error {
	args := m.Called(ctx, productID, quantity)
	return args.Error(0)
}

func (m *MockInventoryRepository) GetProduct(ctx context.Context, productID string) (*repository.Product, error) {
	args := m.Called(ctx, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.Product), args.Error(1)
}

func (m *MockInventoryRepository) CreateProduct(ctx context.Context, product *repository.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockInventoryRepository) CreateInventory(ctx context.Context, inventory *repository.Inventory) error {
	args := m.Called(ctx, inventory)
	return args.Error(0)
}

func TestCheckStock_Success(t *testing.T) {
	// Create mock
	repo := new(MockInventoryRepository)

	// Create service
	inventoryService := service.NewInventoryService(repo)

	// Set up expectations
	repo.On("CheckStock", mock.Anything, "product123", 2).Return(true, nil)

	// Call service
	available, err := inventoryService.CheckStock(context.Background(), "product123", 2)

	// Assert expectations
	assert.NoError(t, err)
	assert.True(t, available)

	// Verify mock
	repo.AssertExpectations(t)
}

func TestCheckStock_Unavailable(t *testing.T) {
	// Create mock
	repo := new(MockInventoryRepository)

	// Create service
	inventoryService := service.NewInventoryService(repo)

	// Set up expectations
	repo.On("CheckStock", mock.Anything, "product123", 2).Return(false, nil)

	// Call service
	available, err := inventoryService.CheckStock(context.Background(), "product123", 2)

	// Assert expectations
	assert.NoError(t, err)
	assert.False(t, available)

	// Verify mock
	repo.AssertExpectations(t)
}

func TestCheckStock_Error(t *testing.T) {
	// Create mock
	repo := new(MockInventoryRepository)

	// Create service
	inventoryService := service.NewInventoryService(repo)

	// Set up expectations
	repo.On("CheckStock", mock.Anything, "product123", 2).Return(false, errors.New("database error"))

	// Call service
	available, err := inventoryService.CheckStock(context.Background(), "product123", 2)

	// Assert expectations
	assert.Error(t, err)
	assert.False(t, available)
	assert.Contains(t, err.Error(), "database error")

	// Verify mock
	repo.AssertExpectations(t)
}

func TestReserveStock_Success(t *testing.T) {
	// Create mock
	repo := new(MockInventoryRepository)

	// Create service
	inventoryService := service.NewInventoryService(repo)

	// Set up expectations
	repo.On("CheckStock", mock.Anything, "product123", 2).Return(true, nil)
	repo.On("ReserveStock", mock.Anything, "product123", 2).Return(nil)

	// Call service
	err := inventoryService.ReserveStock(context.Background(), "product123", 2, "order123")

	// Assert expectations
	assert.NoError(t, err)

	// Verify mock
	repo.AssertExpectations(t)
}

func TestReserveStock_Unavailable(t *testing.T) {
	// Create mock
	repo := new(MockInventoryRepository)

	// Create service
	inventoryService := service.NewInventoryService(repo)

	// Set up expectations
	repo.On("CheckStock", mock.Anything, "product123", 2).Return(false, nil)

	// Call service
	err := inventoryService.ReserveStock(context.Background(), "product123", 2, "order123")

	// Assert expectations
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient stock")

	// Verify mock
	repo.AssertExpectations(t)
}

func TestReserveStock_ReservationError(t *testing.T) {
	// Create mock
	repo := new(MockInventoryRepository)

	// Create service
	inventoryService := service.NewInventoryService(repo)

	// Set up expectations
	repo.On("CheckStock", mock.Anything, "product123", 2).Return(true, nil)
	repo.On("ReserveStock", mock.Anything, "product123", 2).Return(errors.New("reservation error"))

	// Call service
	err := inventoryService.ReserveStock(context.Background(), "product123", 2, "order123")

	// Assert expectations
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reservation error")

	// Verify mock
	repo.AssertExpectations(t)
}

func TestReleaseStock_Success(t *testing.T) {
	// Create mock
	repo := new(MockInventoryRepository)

	// Create service
	inventoryService := service.NewInventoryService(repo)

	// Set up expectations
	repo.On("ReleaseStock", mock.Anything, "product123", 2).Return(nil)

	// Call service
	err := inventoryService.ReleaseStock(context.Background(), "product123", 2, "order123")

	// Assert expectations
	assert.NoError(t, err)

	// Verify mock
	repo.AssertExpectations(t)
}

func TestReleaseStock_Error(t *testing.T) {
	// Create mock
	repo := new(MockInventoryRepository)

	// Create service
	inventoryService := service.NewInventoryService(repo)

	// Set up expectations
	repo.On("ReleaseStock", mock.Anything, "product123", 2).Return(errors.New("release error"))

	// Call service
	err := inventoryService.ReleaseStock(context.Background(), "product123", 2, "order123")

	// Assert expectations
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "release error")

	// Verify mock
	repo.AssertExpectations(t)
}
