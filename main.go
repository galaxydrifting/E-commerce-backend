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

	"e-commerce/di"
	"e-commerce/docs"
	"e-commerce/migrations"
	"e-commerce/routes"

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

	// Initialize dependency injection container
	container, err := di.Initialize(envFile)
	if err != nil {
		log.Fatal("Error initializing container:", err)
	}

	// Run database migrations with the injected DB instance
	migrations.Migrate(container.DB)

	r := gin.Default()

	// Setup routes using the container
	routes.SetupAuthRoutes(r, container.AuthController)

	// Swagger documentation route
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
