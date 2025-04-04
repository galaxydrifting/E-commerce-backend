package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func loadEnv(envFile string) error {
	if envFile == "" {
		envFile = ".env"
	}
	return godotenv.Load(envFile)
}

func ConnectDB(envFile string) {
	err := loadEnv(envFile)
	if err != nil {
		log.Fatal("Error loading env file: ", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	DB = db
	log.Printf("Connected Successfully to Database: %s\n", dbName)
}

func GetDB() *gorm.DB {
	return DB
}
