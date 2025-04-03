// @title           E-Commerce API
// @version         1.0
// @description     This is a sample e-commerce server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

package main

import (
	"e-commerce/configs"
	"e-commerce/migrations"
	"e-commerce/routes"

	_ "e-commerce/docs" // swagger docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Initialize database
	configs.ConnectDB()

	// Run migrations
	migrations.Migrate()

	// Setup router
	r := gin.Default()

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup routes
	routes.SetupAuthRoutes(r)

	// Run the server
	r.Run(":8080")
}
