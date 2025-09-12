package seed

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fardannozami/golang-microservice/inventory-service/repository"
)

// SeedData populates the database with initial data
func SeedData(repo repository.InventoryRepository) error {
	ctx := context.Background()

	// Sample products
	products := []repository.Product{
		{
			ID:          "prod-001",
			Name:        "Laptop Gaming",
			Description: "Laptop gaming dengan spesifikasi tinggi",
			Price:       15000000,
		},
		{
			ID:          "prod-002",
			Name:        "Smartphone",
			Description: "Smartphone dengan kamera 108MP",
			Price:       8000000,
		},
		{
			ID:          "prod-003",
			Name:        "Headphone Bluetooth",
			Description: "Headphone dengan noise cancelling",
			Price:       2000000,
		},
		{
			ID:          "prod-004",
			Name:        "Smart Watch",
			Description: "Jam tangan pintar dengan fitur kesehatan",
			Price:       3500000,
		},
		{
			ID:          "prod-005",
			Name:        "Wireless Earbuds",
			Description: "Earbuds tanpa kabel dengan kualitas suara premium",
			Price:       1800000,
		},
	}

	// Insert products
	for _, product := range products {
		productCopy := product // Create a copy to avoid issues with loop variable capture
		err := repo.CreateProduct(ctx, &productCopy)
		if err != nil {
			log.Printf("Warning: Failed to seed product %s: %v", product.ID, err)
			// Continue with other products even if one fails
			continue
		}
		fmt.Printf("Seeded product: %s - %s\n", product.ID, product.Name)
	}

	// Sample inventory
	inventories := []repository.Inventory{
		{
			ProductID: "prod-001",
			Quantity:  10,
			Reserved:  0,
			UpdatedAt: time.Now(),
		},
		{
			ProductID: "prod-002",
			Quantity:  20,
			Reserved:  0,
			UpdatedAt: time.Now(),
		},
		{
			ProductID: "prod-003",
			Quantity:  30,
			Reserved:  0,
			UpdatedAt: time.Now(),
		},
		{
			ProductID: "prod-004",
			Quantity:  15,
			Reserved:  0,
			UpdatedAt: time.Now(),
		},
		{
			ProductID: "prod-005",
			Quantity:  25,
			Reserved:  0,
			UpdatedAt: time.Now(),
		},
	}

	// Insert inventory
	for _, inventory := range inventories {
		inventoryCopy := inventory // Create a copy to avoid issues with loop variable capture
		err := repo.CreateInventory(ctx, &inventoryCopy)
		if err != nil {
			log.Printf("Warning: Failed to seed inventory for product %s: %v", inventory.ProductID, err)
			// Continue with other inventory items even if one fails
			continue
		}
		fmt.Printf("Seeded inventory for product: %s - Quantity: %d\n", inventory.ProductID, inventory.Quantity)
	}

	return nil
}