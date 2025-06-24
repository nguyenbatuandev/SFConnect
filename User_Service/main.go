package main

import (
    "log"
    "user_service/internal/config"
    "user_service/internal/database"
    "user_service/internal/handler"
    "user_service/internal/middleware"
    "user_service/internal/repository"
    "user_service/internal/service"
    "strconv"
    "github.com/gin-gonic/gin"
)

func main() {
    // Load config
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }

    // Initialize database
    db, err := database.InitDB(cfg)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Initialize repository
    userRepo := repository.NewPostgresRepository(db)

    // Initialize services
    expiresIn,_ := strconv.Atoi(cfg.JwtConfig.ExpiresIn)

    authService := service.NewJWTauthService(cfg.JwtConfig.SecretKey, expiresIn)

    userService := service.NewUserService(userRepo, authService)

    // Initialize handlers
    userHandler := handler.NewUserHandler(userService, authService)

    // Setup Gin
    gin.SetMode(cfg.ServerConfig.GinMode)
    router := gin.Default()

    // Middleware
    router.Use(middleware.CORSMiddleware())


    // Routes
    v1 := router.Group("/api/v1")
    {

        // Public routes
        v1.POST("/register", userHandler.RegisterUser)
        v1.POST("/login", userHandler.LoginUser)
        
        protected := v1.Group("/")
        protected.Use(middleware.AuthMiddleware(authService), middleware.AuthMiddleware(authService))
        {
            protected.GET("/get-profile", userHandler.GetProfile)
            protected.PATCH("/update-profile", userHandler.UpdateProfile)
            protected.DELETE("/delete-profile", userHandler.DeleteUser)
        }

        // Admin routes
        admin := router.Group("api/admin")
        admin.Use(middleware.AuthMiddleware(authService), middleware.AdminOnlyMiddleware())
        {
            admin.POST("/toggle-user-lock/:id", userHandler.ToggleUserLockByAdmin)
            admin.GET("/get-all-user", userHandler.GetAllUsersByAdmin)
        }
        
        buyer := router.Group("api/buyer")
        buyer.Use(middleware.AuthMiddleware(authService), middleware.BuyerOnlyMiddleware())
        {
            buyer.GET("/get-list-partner", userHandler.GetListPartner)
        }
    }

    log.Printf("Server starting on port %s", cfg.ServerConfig.Port)
    if err := router.Run(":" + cfg.ServerConfig.Port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}