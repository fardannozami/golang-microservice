package main

import (
	"testing"

	"github.com/fardannozami/golang-microservice/order-service/docs"
	"github.com/stretchr/testify/assert"
)

func TestSwaggerInfo(t *testing.T) {
	// Verify Swagger info is properly configured
	assert.Equal(t, "Order Service API", docs.SwaggerInfo.Title)
	assert.Equal(t, "1.0", docs.SwaggerInfo.Version)
	assert.Equal(t, "API for managing orders in the microservice architecture", docs.SwaggerInfo.Description)
	assert.Equal(t, "/api/v1", docs.SwaggerInfo.BasePath)
	assert.Equal(t, []string{"http"}, docs.SwaggerInfo.Schemes)
}