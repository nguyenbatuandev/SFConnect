package main

import (
	"Product_Service/internal/config"
	"Product_Service/internal/database"
	"Product_Service/internal/handler"
	"Product_Service/internal/middleware"
	"Product_Service/internal/repository"
	"Product_Service/internal/service"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	productRepo := repository.NewProductRepository(db)
	authService := service.NewJWTauthService(cfg.JWT.SecretKey)

	redisService := service.NewRedisService(cfg.Redis.Host+":"+cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.Database)
	productService := service.NewProductService(productRepo, authService, redisService)

	productHandler := handler.NewProductHandler(productService, authService)

	gin.SetMode(cfg.Server.GinMode)
	r := gin.Default()
	r.Use(gin.Recovery())

	r.Use(middleware.CORSMiddleware())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/product-id/:id", productHandler.GetProductByID)
		v1.GET("/products", productHandler.GetAllProducts)
		v1.GET(("/product-name/:name"), productHandler.GetProductByName)
	}

	v2 := r.Group("/api/admin")
	{
		v2.Use(middleware.AuthMiddleware(authService), middleware.AdminOnlyMiddleware())
		{
			v2.POST("/products", productHandler.CreateProduct)
			v2.PATCH("/products/:id", productHandler.UpdateProduct)
			v2.DELETE("/products/:id", productHandler.DeleteProduct)
		}
	}

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
