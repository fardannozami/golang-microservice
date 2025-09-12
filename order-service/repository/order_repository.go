package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Order represents an order entity
type Order struct {
	ID        string
	UserID    string
	Status    string
	Items     []OrderItem
	CreatedAt time.Time
	UpdatedAt time.Time
}

// OrderItem represents an order item entity
type OrderItem struct {
	ID        string
	OrderID   string
	ProductID string
	Quantity  int
	Price     float64
}

// OrderRepository defines the interface for order repository operations
type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id string) (*Order, error)
	List(ctx context.Context) ([]*Order, error)
	Update(ctx context.Context, order *Order) error
}

// orderRepository implements OrderRepository interface
type orderRepository struct {
	db *sql.DB
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

// Create creates a new order
func (r *orderRepository) Create(ctx context.Context, order *Order) error {
	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Generate a new UUID if not provided
	if order.ID == "" {
		order.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	order.CreatedAt = now
	order.UpdatedAt = now

	// Insert order
	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO orders (id, user_id, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
		order.ID, order.UserID, order.Status, order.CreatedAt, order.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	// Insert order items
	for i := range order.Items {
		// Generate a new UUID for the item if not provided
		if order.Items[i].ID == "" {
			order.Items[i].ID = uuid.New().String()
		}

		// Set order ID
		order.Items[i].OrderID = order.ID

		// Insert order item
		_, err = tx.ExecContext(
			ctx,
			"INSERT INTO order_items (id, order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4, $5)",
			order.Items[i].ID, order.Items[i].OrderID, order.Items[i].ProductID, order.Items[i].Quantity, order.Items[i].Price,
		)
		if err != nil {
			return fmt.Errorf("failed to insert order item: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetByID gets an order by ID
func (r *orderRepository) GetByID(ctx context.Context, id string) (*Order, error) {
	// Query order
	row := r.db.QueryRowContext(
		ctx,
		"SELECT id, user_id, status, created_at, updated_at FROM orders WHERE id = $1",
		id,
	)

	// Scan order
	order := &Order{}
	err := row.Scan(&order.ID, &order.UserID, &order.Status, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found: %s", id)
		}
		return nil, fmt.Errorf("failed to scan order: %w", err)
	}

	// Query order items
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT id, order_id, product_id, quantity, price FROM order_items WHERE order_id = $1",
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query order items: %w", err)
	}
	defer rows.Close()

	// Scan order items
	for rows.Next() {
		item := OrderItem{}
		err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}
		order.Items = append(order.Items, item)
	}

	return order, nil
}

// List lists all orders
func (r *orderRepository) List(ctx context.Context) ([]*Order, error) {
	// Query orders
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT id, user_id, status, created_at, updated_at FROM orders ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	// Scan orders
	orders := []*Order{}
	for rows.Next() {
		order := &Order{}
		err := rows.Scan(&order.ID, &order.UserID, &order.Status, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	// Query order items for each order
	for _, order := range orders {
		// Query order items
		rows, err := r.db.QueryContext(
			ctx,
			"SELECT id, order_id, product_id, quantity, price FROM order_items WHERE order_id = $1",
			order.ID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to query order items: %w", err)
		}
		defer rows.Close()

		// Scan order items
		for rows.Next() {
			item := OrderItem{}
			err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price)
			if err != nil {
				return nil, fmt.Errorf("failed to scan order item: %w", err)
			}
			order.Items = append(order.Items, item)
		}
	}

	return orders, nil
}

// Update updates an order
func (r *orderRepository) Update(ctx context.Context, order *Order) error {
	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Set updated timestamp
	order.UpdatedAt = time.Now()

	// Update order
	_, err = tx.ExecContext(
		ctx,
		"UPDATE orders SET user_id = $1, status = $2, updated_at = $3 WHERE id = $4",
		order.UserID, order.Status, order.UpdatedAt, order.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}