package database

import (
	"fmt"
	"log"
	"os"
	"task-vix-btpns/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

func ConnectDB() *gorm.DB {
	godotenv.Load(".env")

	DB_HOST := os.Getenv("DB_HOST")
	DB_USER := os.Getenv("DB_USER")
	DB_DRIVER := os.Getenv("DB_DRIVER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")
	DB_PORT := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
	db, err := gorm.Open("mysql", dsn) //Connecting to database

	if err != nil {
		fmt.Printf("Cannot connect to %s database", DB_DRIVER)
		log.Fatal(err)
	}

	err = db.Debug().AutoMigrate(&models.User{}, &models.Photo{}).Error //Migrate the tables to database
	if err != nil {
		log.Fatalf("Migrating table error: %v", err)
	}

	err = db.Debug().Model(&models.Photo{}).AddForeignKey("user_id", "users(id)", "cascade", "cascade").Error //Add foreign key to table
	if err != nil {
		log.Fatalf("Error while attaching foreign key: %v", err)
	}

	return db
}
