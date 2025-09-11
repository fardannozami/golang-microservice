package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/fardannozami/golang-microservice/inventory-service/config"
	inventorypb "github.com/fardannozami/golang-microservice/inventory-service/proto"
	"github.com/fardannozami/golang-microservice/inventory-service/repository"
	"github.com/fardannozami/golang-microservice/inventory-service/server"
	"github.com/fardannozami/golang-microservice/inventory-service/service"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	fmt.Println(cfg)

	// Initialize database connection
	db, err := repository.NewPostgresConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	inventoryRepo := repository.NewInventoryRepository(db)

	// Initialize services
	inventoryService := service.NewInventoryService(inventoryRepo)

	// Initialize gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.ServerPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Register inventory service
	inventoryServer := server.NewInventoryServer(inventoryService)
	inventorypb.RegisterInventoryServiceServer(grpcServer, inventoryServer)

	// Start gRPC server in a goroutine
	go func() {
		log.Printf("Starting inventory service on port %d", cfg.ServerPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Gracefully stop the gRPC server
	log.Println("Shutting down server...")
	grpcServer.GracefulStop()
	log.Println("Server exited properly")
}
