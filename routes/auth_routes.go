package routes

import (
	"e-commerce/controllers"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, authController *controllers.AuthController) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}
}
