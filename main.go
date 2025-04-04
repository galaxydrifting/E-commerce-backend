// @title           E-commerce API
// @version         1.0
// @description     A simple e-commerce API.
// @termsOfService  http://swagger.io/terms/

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"log"
	"os"

	"e-commerce/configs"
	"e-commerce/controllers"
	"e-commerce/docs"
	"e-commerce/repository"
	"e-commerce/routes"
	"e-commerce/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Load environment variables
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}
	if err := godotenv.Load(envFile); err != nil {
		log.Fatal("Error loading .env file")
	}

	configs.ConnectDB(envFile)

	r := gin.Default()

	// Initialize repositories and services
	userRepo := repository.NewGormUserRepository(configs.GetDB())
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	// Setup routes
	routes.SetupAuthRoutes(r, authController)

	// Swagger documentation route
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
