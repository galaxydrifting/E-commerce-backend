package migrations

import (
	"log"

	"e-commerce/configs"
	"e-commerce/models"
)

func Migrate() {
	err := configs.DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}
	log.Println("Database Migration Completed!")
}
