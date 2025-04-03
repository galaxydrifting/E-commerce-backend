package main

import (
	"e-commerce/configs"
	"e-commerce/migrations"
	"e-commerce/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	configs.ConnectDB()

	// Run migrations
	migrations.Migrate()

	// Setup router
	r := gin.Default()

	// Setup routes
	routes.SetupAuthRoutes(r)

	// Run the server
	r.Run(":8080")
}
