package main

import (
	"log"

	"paddle-traceability/blockchain"
	"paddle-traceability/config"
	"paddle-traceability/database"
	"paddle-traceability/handlers"
	"paddle-traceability/middleware"
	"paddle-traceability/services"
	"paddle-traceability/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg := config.Load()

	// Init snowflake
	if err := utils.InitSnowflake(); err != nil {
		log.Fatalf("snowflake init failed: %v", err)
	}

	// Init database
	database.InitDB(&cfg.DB)

	// Init blockchain client
	chainClient, err := blockchain.NewChainClient(&cfg.Blockchain)
	if err != nil {
		log.Fatalf("blockchain client init failed: %v", err)
	}

	// Init service layer
	authService := services.NewAuthService(&cfg.JWT)
	productService := services.NewProductService(chainClient)
	logisticsService := services.NewLogisticsService(chainClient)
	traceService := services.NewTraceService(chainClient)

	// Init handlers
	authHandler := handlers.NewAuthHandler(authService)
	productHandler := handlers.NewProductHandler(productService, traceService)
	logisticsHandler := handlers.NewLogisticsHandler(logisticsService)
	verifyHandler := handlers.NewVerifyHandler(traceService)

	// Init Gin router
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORS())

	// Routes
	api := r.Group("/api/v1")
	{
		// Auth (no token required)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/profile", middleware.JWTAuth(cfg.JWT.Secret), authHandler.GetProfile)
		}

		// Products (token required)
		products := api.Group("/products")
		products.Use(middleware.JWTAuth(cfg.JWT.Secret))
		{
			products.POST("", productHandler.CreateProduct)
			products.GET("", productHandler.ListProducts)
			products.GET("/:id", productHandler.GetProduct)
			products.PUT("/:id", productHandler.UpdateProduct)
			products.GET("/:id/trace", productHandler.GetTrace)
		}

		// Logistics (token required)
		logistics := api.Group("/logistics")
		logistics.Use(middleware.JWTAuth(cfg.JWT.Secret))
		{
			logistics.POST("", logisticsHandler.AddRecord)
			logistics.GET("", logisticsHandler.GetRecords)
		}

		// Anti-counterfeiting verification (public)
		verify := api.Group("/verify")
		{
			verify.GET("/:product_uid", verifyHandler.VerifyProduct)
		}
	}

	log.Println("server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server start failed: %v", err)
	}
}
