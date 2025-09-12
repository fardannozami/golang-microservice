package server

import (
	"context"
	"log"

	inventorypb "github.com/fardannozami/golang-microservice/inventory-service/proto"
	"github.com/fardannozami/golang-microservice/inventory-service/service"
)

// InventoryServer implements the gRPC server for inventory service
type InventoryServer struct {
	inventorypb.UnimplementedInventoryServiceServer
	service service.InventoryService
}

// NewInventoryServer creates a new inventory server
func NewInventoryServer(service service.InventoryService) *InventoryServer {
	return &InventoryServer{service: service}
}

// CheckStock checks if a product is available in inventory
func (s *InventoryServer) CheckStock(ctx context.Context, req *inventorypb.CheckStockRequest) (*inventorypb.CheckStockResponse, error) {
	log.Printf("[inventory-service] CheckStock product_id=%s qty=%d", req.ProductId, req.Quantity)
	// Call service
	available, err := s.service.CheckStock(ctx, req.ProductId, int(req.Quantity))
	if err != nil {
		return &inventorypb.CheckStockResponse{
			Available: false,
			Message:   err.Error(),
		}, nil
	}

	// Return response
	return &inventorypb.CheckStockResponse{
		Available: available,
		Message:   "",
	}, nil
}

// ReserveStock reserves stock for an order
func (s *InventoryServer) ReserveStock(ctx context.Context, req *inventorypb.ReserveStockRequest) (*inventorypb.ReserveStockResponse, error) {
	log.Printf("[inventory-service] ReserveStock product_id=%s qty=%d order_id=%s", req.ProductId, req.Quantity, req.OrderId)
	// Call service
	err := s.service.ReserveStock(ctx, req.ProductId, int(req.Quantity), req.OrderId)
	if err != nil {
		return &inventorypb.ReserveStockResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Return response
	return &inventorypb.ReserveStockResponse{
		Success: true,
		Message: "",
	}, nil
}

// ReleaseStock releases reserved stock
func (s *InventoryServer) ReleaseStock(ctx context.Context, req *inventorypb.ReleaseStockRequest) (*inventorypb.ReleaseStockResponse, error) {
	log.Printf("[inventory-service] ReleaseStock product_id=%s qty=%d order_id=%s", req.ProductId, req.Quantity, req.OrderId)
	// Call service
	err := s.service.ReleaseStock(ctx, req.ProductId, int(req.Quantity), req.OrderId)
	if err != nil {
		return &inventorypb.ReleaseStockResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Return response
	return &inventorypb.ReleaseStockResponse{
		Success: true,
		Message: "",
	}, nil
}
