package db

import (
	"github.com/Phase-R/Phase-R-Backend/db/database"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func Init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DB = database.InitDB()
}
