package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fardannozami/golang-microservice/order-service/config"
	"github.com/fardannozami/golang-microservice/order-service/docs"
	"github.com/fardannozami/golang-microservice/order-service/handler"
	"github.com/fardannozami/golang-microservice/order-service/repository"
	"github.com/fardannozami/golang-microservice/order-service/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	fmt.Println(cfg)
	
	// Swagger configuration
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", cfg.ServerPort)
	docs.SwaggerInfo.BasePath = "/api/v1"

	// Initialize database connection
	db, err := repository.NewPostgresConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	orderRepo := repository.NewOrderRepository(db)

	// Initialize inventory client
	inventoryClient, err := service.NewInventoryClient(cfg.InventoryServiceURL)
	if err != nil {
		log.Fatalf("Failed to create inventory client: %v", err)
	}
	defer inventoryClient.Close()

	// Initialize services
	orderService := service.NewOrderService(orderRepo, inventoryClient)

	// Initialize handlers
	orderHandler := handler.NewOrderHandler(orderService)

	// Initialize router
	router := gin.Default()

	// Register middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Register routes
	v1 := router.Group("/api/v1")
	{
		orders := v1.Group("/orders")
		{
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("", orderHandler.ListOrders)
			orders.GET("/:id", orderHandler.GetOrder)
		}
	}
	
	// Swagger documentation route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: router,
	}

	// Run server in a goroutine
	go func() {
		log.Printf("Starting order service on port %d", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
