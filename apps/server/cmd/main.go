package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/price-comparison/server/internal/handler"
	"github.com/price-comparison/server/internal/repository"
)

func main() {
	// Database configuration
	dbConfig := repository.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     5432,
		User:     getEnv("DB_USER", "admin"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "price_comparison"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Initialize database
	db, err := repository.NewDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	storeRepo := repository.NewStoreRepository(db)
	productRepo := repository.NewProductRepository(db)
	priceRepo := repository.NewPriceRepository(db)

	// Initialize handlers
	storeHandler := handler.NewStoreHandler(storeRepo)
	productHandler := handler.NewProductHandler(productRepo, priceRepo)

	// Setup Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := r.Group("/api")
	{
		// Store routes
		stores := api.Group("/stores")
		{
			stores.GET("", storeHandler.GetAllStores)
			stores.GET("/nearby", storeHandler.GetNearbyStores)
			stores.GET("/:id", storeHandler.GetStoreByID)
		}

		// Product routes
		products := api.Group("/products")
		{
			products.GET("", productHandler.GetAllProducts)
			products.GET("/search", productHandler.SearchProducts)
			products.GET("/:id", productHandler.GetProductByID)
			products.GET("/:id/prices", productHandler.GetProductPrices)
		}
	}

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
