package main

import (
	"os"
	"task-vix-btpns/database"
	"task-vix-btpns/models"
	"task-vix-btpns/router"
)

func main() {
	db := database.ConnectDB()
	db.AutoMigrate(&models.User{})

	r := router.InitRoutes(db)
	r.Run(":" + os.Getenv("PORT"))
}
