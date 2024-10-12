package database

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/Phase-R/Phase-R-Backend/db/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBLock sync.Mutex

func InitDB() *gorm.DB {
	DBLock.Lock()
	defer DBLock.Unlock()
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := 5432

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=require", host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	db.AutoMigrate(&models.User{}, &models.Activities{}, &models.ActivityType{}, &models.Drill{}, &models.Thoughts{})
	return db
}