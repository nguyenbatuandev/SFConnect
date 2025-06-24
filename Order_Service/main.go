package main

import (
	"Order_Service/internal/config"
	"Order_Service/internal/database"
	"Order_Service/internal/handle"
	"Order_Service/internal/middleware"
	"Order_Service/internal/repository"
	"Order_Service/internal/service"
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
	callUserService := service.NewCallService()

	//repo
	orderRepository := repository.NewOrderRepository(db, callUserService)
	partnerRepository := repository.NewPartnerRepository(db, callUserService)
	adminRepository := repository.NewAdminRepository(db, callUserService)
	buyerRepository := repository.NewBuyerRepository(db, callUserService)
	//service
	authService := service.NewJWTauthService(cfg.JWT.SecretKey)
	partnerCommis := repository.NewPartnerCommissionRepository(db, orderRepository)
	partnerCommisService := service.NewPartnerCommissionService(partnerCommis)
	buyerService := service.NewBuyerService(orderRepository, partnerCommisService, buyerRepository)
	partnerService := service.NewPartnerService(orderRepository, partnerRepository)
	adminService := service.NewAdminService(orderRepository, partnerCommisService, adminRepository)

	buyerHandle := handle.NewBuyerHandler(buyerService , callUserService)
	partnerHandle := handle.NewPartnerHandler(partnerService, partnerCommisService)
	adminHandle := handle.NewAdminHandler(adminService, partnerCommisService)
	// Initialize HTTP server
	gin.SetMode(cfg.Server.GinMode)
	r := gin.Default()
	r.Use(gin.Recovery())

	r.Use(middleware.CORSMiddleware())

	// Buyer routes
	buyerGroup := r.Group("api/buyer")
	{

		buyerGroup.Use(middleware.AuthMiddleware(authService), middleware.BuyerOnlyMiddleware())
		{
			buyerGroup.GET("/orders/:status", buyerHandle.GetAllOrdersByBuyerID)
			buyerGroup.POST("/orders", buyerHandle.CreateOrder)
			buyerGroup.PUT("/orders/:orderItemId", buyerHandle.UpdateOrderStatus)
			buyerGroup.GET("/partners", buyerHandle.GetListPartner)
		}
	}
	// Partner routes
	partnerGroup := r.Group("api/partner")
	{
		partnerGroup.Use(middleware.AuthMiddleware(authService), middleware.PartnerOnlyMiddleware())
		{
			partnerGroup.GET("/orders/:status", partnerHandle.GetAllOrdersByPartnerID)
			partnerGroup.GET("/commission", partnerHandle.GetCommissionByPartnerID)
			partnerGroup.PUT("/orders/:orderItemId", partnerHandle.UpdateOrderStatus)
		}
	}

	// Admin routes
	adminGroup := r.Group("api/admin")
	{
		adminGroup.Use(middleware.AuthMiddleware(authService), middleware.AdminOnlyMiddleware())
		{
			adminGroup.GET("/all-orders", adminHandle.GetAllOrders)
			adminGroup.GET("/orders/:orderId", adminHandle.GetOrderByID)
			adminGroup.GET("/commission/:orderItemId", adminHandle.GetCommissionByOrderItemID)
			adminGroup.GET("/commission-pnid/:patrnerId", adminHandle.GetCommissionByPartnerID)
		}
	}

	log.Printf("Server starting on port %s", cfg.Server.Port)
    if err := r.Run(":" + cfg.Server.Port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
