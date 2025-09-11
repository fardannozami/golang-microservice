package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Product represents a product entity
type Product struct {
	ID          string
	Name        string
	Description string
	Price       float64
}

// Inventory represents an inventory entity
type Inventory struct {
	ProductID string
	Quantity  int
	Reserved  int
	UpdatedAt time.Time
}

// InventoryRepository defines the interface for inventory repository operations
type InventoryRepository interface {
	CheckStock(ctx context.Context, productID string, quantity int) (bool, error)
	ReserveStock(ctx context.Context, productID string, quantity int) error
	ReleaseStock(ctx context.Context, productID string, quantity int) error
	GetProduct(ctx context.Context, productID string) (*Product, error)
	CreateProduct(ctx context.Context, product *Product) error
	CreateInventory(ctx context.Context, inventory *Inventory) error
}

// inventoryRepository implements InventoryRepository interface
type inventoryRepository struct {
	db *sql.DB
}

// NewInventoryRepository creates a new inventory repository
func NewInventoryRepository(db *sql.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

// CheckStock checks if a product is available in inventory
func (r *inventoryRepository) CheckStock(ctx context.Context, productID string, quantity int) (bool, error) {
	// Query inventory
	row := r.db.QueryRowContext(
		ctx,
		"SELECT quantity, reserved FROM inventory WHERE product_id = $1",
		productID,
	)

	// Scan inventory
	var inventoryQuantity, reserved int
	err := row.Scan(&inventoryQuantity, &reserved)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to scan inventory: %w", err)
	}

	// Check if available quantity is sufficient
	available := inventoryQuantity - reserved
	return available >= quantity, nil
}

// ReserveStock reserves stock for an order
func (r *inventoryRepository) ReserveStock(ctx context.Context, productID string, quantity int) error {
	// Start a transaction
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Query inventory with lock
	row := tx.QueryRowContext(
		ctx,
		"SELECT quantity, reserved FROM inventory WHERE product_id = $1 FOR UPDATE",
		productID,
	)

	// Scan inventory
	var inventoryQuantity, reserved int
	err = row.Scan(&inventoryQuantity, &reserved)
	if err != nil {
		return fmt.Errorf("failed to scan inventory: %w", err)
	}

	// Check if available quantity is sufficient
	available := inventoryQuantity - reserved
	if available < quantity {
		return fmt.Errorf("insufficient stock: available %d, requested %d", available, quantity)
	}

	// Update reserved quantity
	_, err = tx.ExecContext(
		ctx,
		"UPDATE inventory SET reserved = reserved + $1, updated_at = $2 WHERE product_id = $3",
		quantity, time.Now(), productID,
	)
	if err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ReleaseStock releases reserved stock
func (r *inventoryRepository) ReleaseStock(ctx context.Context, productID string, quantity int) error {
	// Start a transaction
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Query inventory with lock
	row := tx.QueryRowContext(
		ctx,
		"SELECT reserved FROM inventory WHERE product_id = $1 FOR UPDATE",
		productID,
	)

	// Scan inventory
	var reserved int
	err = row.Scan(&reserved)
	if err != nil {
		return fmt.Errorf("failed to scan inventory: %w", err)
	}

	// Check if reserved quantity is sufficient
	if reserved < quantity {
		return fmt.Errorf("invalid release: reserved %d, requested %d", reserved, quantity)
	}

	// Update reserved quantity
	_, err = tx.ExecContext(
		ctx,
		"UPDATE inventory SET reserved = reserved - $1, updated_at = $2 WHERE product_id = $3",
		quantity, time.Now(), productID,
	)
	if err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetProduct gets a product by ID
func (r *inventoryRepository) GetProduct(ctx context.Context, productID string) (*Product, error) {
	// Query product
	row := r.db.QueryRowContext(
		ctx,
		"SELECT id, name, description, price FROM products WHERE id = $1",
		productID,
	)

	// Scan product
	product := &Product{}
	err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found: %s", productID)
		}
		return nil, fmt.Errorf("failed to scan product: %w", err)
	}

	return product, nil
}

// CreateProduct creates a new product
func (r *inventoryRepository) CreateProduct(ctx context.Context, product *Product) error {
	// Insert product
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO products (id, name, description, price) VALUES ($1, $2, $3, $4)",
		product.ID, product.Name, product.Description, product.Price,
	)
	if err != nil {
		return fmt.Errorf("failed to insert product: %w", err)
	}

	return nil
}

// CreateInventory creates a new inventory entry
func (r *inventoryRepository) CreateInventory(ctx context.Context, inventory *Inventory) error {
	// Set timestamp
	inventory.UpdatedAt = time.Now()

	// Insert inventory
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO inventory (product_id, quantity, reserved, updated_at) VALUES ($1, $2, $3, $4)",
		inventory.ProductID, inventory.Quantity, inventory.Reserved, inventory.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert inventory: %w", err)
	}

	return nil
}
