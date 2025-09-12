package main

import (
	"fmt"
	"log"

	"github.com/fardannozami/golang-microservice/inventory-service/config"
	"github.com/fardannozami/golang-microservice/inventory-service/repository"
	"github.com/fardannozami/golang-microservice/inventory-service/seed"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := repository.NewPostgresConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repository
	inventoryRepo := repository.NewInventoryRepository(db)

	// Run seeder
	fmt.Println("Starting to seed inventory data...")
	if err := seed.SeedData(inventoryRepo); err != nil {
		log.Fatalf("Failed to seed data: %v", err)
	}

	fmt.Println("Seeding completed successfully!")
}