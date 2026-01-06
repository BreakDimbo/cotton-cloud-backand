package api

import (
	"cotton-cloud-backend/internal/api/handlers"
	"cotton-cloud-backend/internal/api/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NewRouter creates and configures the Gin router
func NewRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// Middleware
	router.Use(corsMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Public authentication routes (no auth required)
		auth := v1.Group("/auth")
		{
			authHandler := handlers.NewAuthHandler(db)
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected routes (auth middleware with demo fallback)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Clothing routes
			clothing := protected.Group("/clothing")
			{
				clothingHandler := handlers.NewClothingHandler(db)
				clothing.GET("", clothingHandler.List)
				clothing.POST("", clothingHandler.Create)
				clothing.GET("/:id", clothingHandler.Get)
				clothing.PUT("/:id", clothingHandler.Update)
				clothing.DELETE("/:id", clothingHandler.Delete)
				clothing.POST("/:id/wash", clothingHandler.Wash)
				clothing.POST("/:id/wear", clothingHandler.IncrementWear)
			}

			// Avatar routes
			avatars := protected.Group("/avatars")
			{
				avatarHandler := handlers.NewAvatarHandler(db)
				avatars.GET("", avatarHandler.List)
				avatars.POST("", avatarHandler.Create)
				avatars.GET("/:id", avatarHandler.Get)
				avatars.PUT("/:id", avatarHandler.Update)
				avatars.DELETE("/:id", avatarHandler.Delete)
				avatars.POST("/:id/activate", avatarHandler.Activate)
			}

			// Outfit routes
			outfits := protected.Group("/outfits")
			{
				outfitHandler := handlers.NewOutfitHandler(db)
				outfits.GET("", outfitHandler.List)
				outfits.POST("", outfitHandler.Create)
				outfits.GET("/:date", outfitHandler.GetByDate)
				outfits.PUT("/:id", outfitHandler.Update)
				outfits.DELETE("/:id", outfitHandler.Delete)
			}

			// AI proxy routes
			ai := protected.Group("/ai")
			{
				aiHandler := handlers.NewAIHandler()
				ai.POST("/analyze", aiHandler.AnalyzeClothing)
				ai.POST("/cutout", aiHandler.GenerateCutout)
				ai.POST("/avatar", aiHandler.GenerateAvatar)
				ai.POST("/collage", aiHandler.GenerateCollage)
				ai.POST("/tryon", aiHandler.VirtualTryOn)
			}
		}
	}

	return router
}

// corsMiddleware handles CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
