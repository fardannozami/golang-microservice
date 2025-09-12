package handler

import (
	"net/http"

	"github.com/fardannozami/golang-microservice/order-service/service"
	"github.com/gin-gonic/gin"
)

// @title Order Service API
// @version 1.0
// @description API for managing orders in the microservice architecture
// @host localhost:8080
// @BasePath /api/v1

// OrderHandler handles HTTP requests for orders
type OrderHandler struct {
	orderService service.OrderService
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// CreateOrderRequest represents a request to create an order
type CreateOrderRequest struct {
	UserID string                   `json:"user_id" binding:"required" example:"123e4567-e89b-12d3-a456-426655440000"`
	Items  []CreateOrderItemRequest `json:"items" binding:"required,dive"`
}

// CreateOrderItemRequest represents a request to create an order item
type CreateOrderItemRequest struct {
	ProductID string  `json:"product_id" binding:"required" example:"prod-001"`
	Quantity  int     `json:"quantity" binding:"required,gt=0" example:"2"`
	Price     float64 `json:"price" binding:"required,gt=0" example:"10.99"`
}

// OrderResponse represents an order response
type OrderResponse struct {
	ID        string              `json:"id"`
	UserID    string              `json:"user_id"`
	Status    string              `json:"status"`
	Items     []OrderItemResponse `json:"items"`
	CreatedAt string              `json:"created_at"`
	UpdatedAt string              `json:"updated_at"`
}

// OrderItemResponse represents an order item response
type OrderItemResponse struct {
	ID        string  `json:"id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order with the provided items
// @Tags orders
// @Accept json
// @Produce json
// @Param order body CreateOrderRequest true "Order details"
// @Success 201 {object} OrderResponse
// @Failure 400 {object} map[string]interface{}
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	// Parse request
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert request to service request
	serviceReq := &service.CreateOrderRequest{
		UserID: req.UserID,
		Items:  make([]service.OrderItemRequest, len(req.Items)),
	}

	// Convert order items
	for i, item := range req.Items {
		serviceReq.Items[i] = service.OrderItemRequest{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	// Create order
	order, err := h.orderService.CreateOrder(c.Request.Context(), serviceReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert order to response
	resp := OrderResponse{
		ID:        order.ID,
		UserID:    order.UserID,
		Status:    order.Status,
		Items:     make([]OrderItemResponse, len(order.Items)),
		CreatedAt: order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: order.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Convert order items
	for i, item := range order.Items {
		resp.Items[i] = OrderItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	c.JSON(http.StatusCreated, resp)
}

// GetOrder godoc
// @Summary Get an order by ID
// @Description Get an order by its ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} OrderResponse
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrder(c *gin.Context) {
	// Get order ID from path
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order ID is required"})
		return
	}

	// Get order
	order, err := h.orderService.GetOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Convert order to response
	resp := OrderResponse{
		ID:        order.ID,
		UserID:    order.UserID,
		Status:    order.Status,
		Items:     make([]OrderItemResponse, len(order.Items)),
		CreatedAt: order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: order.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Convert order items
	for i, item := range order.Items {
		resp.Items[i] = OrderItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	c.JSON(http.StatusOK, resp)
}

// ListOrders godoc
// @Summary List all orders
// @Description Get a list of all orders
// @Tags orders
// @Accept json
// @Produce json
// @Success 200 {array} OrderResponse
// @Router /orders [get]
func (h *OrderHandler) ListOrders(c *gin.Context) {
	// List orders
	orders, err := h.orderService.ListOrders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert orders to response
	resp := make([]OrderResponse, len(orders))
	for i, order := range orders {
		resp[i] = OrderResponse{
			ID:        order.ID,
			UserID:    order.UserID,
			Status:    order.Status,
			Items:     make([]OrderItemResponse, len(order.Items)),
			CreatedAt: order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: order.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		// Convert order items
		for j, item := range order.Items {
			resp[i].Items[j] = OrderItemResponse{
				ID:        item.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Price,
			}
		}
	}

	c.JSON(http.StatusOK, resp)
}
