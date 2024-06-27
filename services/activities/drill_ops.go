package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Phase-R/Phase-R-Backend/db/models"
	"github.com/joho/godotenv"

	// "github.com/Phase-R/Phase-R-Backend/services/activities/configs"
	"github.com/gin-gonic/gin"
	"github.com/nrednav/cuid2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// type Activity interface {
// CreateActivity(ctx *gofr.Context, activity *models.Activity) (*models.Activity, error)
// GetActivity(ctx *gofr.Context, UUID string) (*models.Activity, error)
// GetActivitiesByType(ctx *gofr.Context) ([]models.Activity, error)
// }

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func InitDB() *gorm.DB {
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
	return db
}

func CreateActivity(ctx *gin.Context) {
	const uniqueViolation = "23505"

	var drill models.Drill
	if err := ctx.ShouldBindJSON(&drill); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := cuid2.Generate()
	if id == "" {
		return
	}

	drill.ID = id

	res := db.Create(&drill)
	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Drill created successfully."})
}

func GetActivityByType(ctx *gin.Context) {
	type_act := ctx.Param("type")
	var drills []models.Drill
	res := db.Raw("SELECT * FROM drills WHERE type_act = ?", type_act).Scan(&drills)
	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}

	if len(drills) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No drills found!"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": drills})
	return
}

func DeleteActivity(ctx *gin.Context) {
	id := ctx.Param("id")
	var drill models.Drill
	res := db.Where("id = ?", id).First(&drill)

	if res.Error != nil {
		errorMsg := fmt.Sprintf("No drills with id: %s were found!", id)
		ctx.JSON(http.StatusNotFound, gin.H{"message": errorMsg})
		return
	}

	del := db.Delete(&drill)
	if del.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": del.Error.Error()})
		return
	}

	successMsg := fmt.Sprintf("Deleted drill: %s successfully!", drill)
	ctx.JSON(http.StatusOK, gin.H{"message": successMsg})
}

func GetActivity(ctx *gin.Context) {
	id := ctx.Param("id")
	var drill models.Drill
	res := db.Raw("SELECT * FROM drills WHERE id = ?", id).Scan(&drill)
	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": drill})
}

func UpdateActivity(ctx *gin.Context) {
	id := ctx.Param("id")
	var drill models.Drill
	var changes models.Drill

	err := ctx.ShouldBindJSON(&changes)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := db.Where("id = ?", id).First(&drill)

	if res.Error != nil {
		notFoundMsg := fmt.Sprintf("No drills with the id: %s were found!", id)
		ctx.JSON(http.StatusNotFound, gin.H{"message": notFoundMsg})
		return
	}

	update := res.Updates(changes)
	if update.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": update.Error.Error()})
		return
	}

	successMsg := fmt.Sprintf("Updated drill with id: %s successfully!", id)
	ctx.JSON(http.StatusOK, gin.H{"message": successMsg})
}

func main() {
	r := gin.Default()
	db = InitDB()
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	r.POST("/create_drill", CreateActivity)
	r.GET("/get_drill/:id", GetActivity)
	r.GET("/get_drills_by_type/:type", GetActivityByType)
	r.PATCH("/update_drill/:id", UpdateActivity)
	r.DELETE("/delete_drill/:id", DeleteActivity)
	r.Run()
}
