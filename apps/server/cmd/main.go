package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/price-comparison/server/internal/cache"
	"github.com/price-comparison/server/internal/config"
	"github.com/price-comparison/server/internal/handler"
	"github.com/price-comparison/server/internal/logger"
	"github.com/price-comparison/server/internal/metrics"
	"github.com/price-comparison/server/internal/middleware"
	"github.com/price-comparison/server/internal/repository"
	"github.com/price-comparison/server/internal/usecase"
)

func main() {
	cfg := config.Load()

	// Initialize database
	db, err := repository.NewDatabase(repository.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.DBName,
		SSLMode:  cfg.DB.SSLMode,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	storeRepo := repository.NewStoreRepository(db)
	productRepo := repository.NewProductRepository(db)
	priceRepo := repository.NewPriceRepository(db)

	var cacheAdapter usecase.Cache
	redisClient, err := cache.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Printf("Redis unavailable, caching disabled: %v", err)
	} else {
		cacheAdapter = cache.NewRedisCache(redisClient)
	}
	cacheTTL := time.Duration(cfg.Cache.TTLSeconds) * time.Second

	// Initialize usecases
	storeUsecase := usecase.NewStoreUsecase(storeRepo, cacheAdapter, cacheTTL)
	productUsecase := usecase.NewProductUsecase(productRepo, cacheAdapter, cacheTTL)
	priceUsecase := usecase.NewPriceUsecase(priceRepo)

	// Initialize handlers
	storeHandler := handler.NewStoreHandler(storeUsecase, priceUsecase)
	productHandler := handler.NewProductHandler(productUsecase, priceUsecase)

	// Setup Gin router
	appLogger := logger.New(cfg.Log.Level)
	metrics.Init()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logging(appLogger))
	r.Use(metrics.Middleware())

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Server.CORSOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Metrics
	r.GET(cfg.Server.MetricsRoute, metrics.Handler())

	// API routes
	api := r.Group("/api")
	api.Use(middleware.APIKeyAuth(cfg.Auth.APIKey))
	{
		// Store routes
		stores := api.Group("/stores")
		{
			stores.GET("", storeHandler.GetAllStores)
			stores.GET("/nearby", storeHandler.GetNearbyStores)
			stores.GET("/:id", storeHandler.GetStoreByID)
			stores.GET("/:id/price-stats", storeHandler.GetStorePriceStats)
			stores.GET("/:id/prices", storeHandler.GetStorePrices)
		}

		// Product routes
		products := api.Group("/products")
		{
			products.GET("", productHandler.GetAllProducts)
			products.GET("/categories", productHandler.GetCategories)
			products.GET("/search", productHandler.SearchProducts)
			products.GET("/:id", productHandler.GetProductByID)
			products.GET("/:id/prices", productHandler.GetProductPrices)
		}
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
