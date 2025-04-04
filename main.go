// @title           E-Commerce API
// @version         1.0
// @description     A simple e-commerce API.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter the token with the `Bearer: ` prefix, e.g. "Bearer abcde12345".

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

	// Setup routes using the container and middleware
	routes.SetupAuthRoutes(r, container.AuthController, container.AuthMiddleware)

	// Swagger documentation route
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
