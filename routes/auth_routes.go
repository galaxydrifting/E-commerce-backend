package routes

import (
	"e-commerce/controllers"
	"e-commerce/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, authController *controllers.AuthController, authMiddleware *middlewares.AuthMiddleware) {
	// 修改為 v1 版本的 API 路徑
	v1 := router.Group("/api/v1")
	auth := v1.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)

		// Protected routes
		protected := auth.Group("")
		protected.Use(authMiddleware.Handle())
		{
			protected.POST("/logout", authController.Logout)
			protected.GET("/profile", authController.GetProfile)
			protected.PUT("/profile", authController.UpdateProfile)
		}
	}
}
