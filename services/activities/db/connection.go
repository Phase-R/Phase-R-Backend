package db

import (
	"log"

	"github.com/Phase-R/Phase-R-Backend/db/database"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	DB = database.InitDB()
}
