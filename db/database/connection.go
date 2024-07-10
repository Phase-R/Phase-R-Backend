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

	// Auto migrate models
	err = db.AutoMigrate(&models.User{}, &models.Activities{}, &models.ActivityType{}, &models.Drill{})
	if err != nil {
		log.Fatalf("failed to auto migrate models: %v", err)
	}

	// Create trigger function to delete expired OTPs
	createTriggerFunction := `
	CREATE OR REPLACE FUNCTION delete_expired_otps() RETURNS trigger AS $$
	BEGIN
	    IF NEW.otp IS NOT NULL AND NEW.updated_at IS NOT NULL THEN
	        IF (NEW.updated_at + INTERVAL '5 minutes') < NOW() THEN
	            NEW.otp := NULL;
	        END IF;
	    END IF;
	    RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;
	`
	db.Exec(createTriggerFunction)

	// Create trigger for delete_expired_otps function
	createTrigger := `
	CREATE TRIGGER trigger_delete_expired_otps
	BEFORE INSERT OR UPDATE ON users
	FOR EACH ROW
	EXECUTE FUNCTION delete_expired_otps();
	`
	db.Exec(createTrigger)

	return db
}
