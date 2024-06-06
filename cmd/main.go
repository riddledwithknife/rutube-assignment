package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"rutube-assignment/internal/config"
	"rutube-assignment/internal/models"
	"rutube-assignment/internal/routes"
	"rutube-assignment/internal/services"
)

func main() {
	db, err := gorm.Open(postgres.Open(config.GetDBConnectionString()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db.AutoMigrate(&models.User{}, &models.Subscription{})

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	routes.SetupRoutes(r, db)

	emailSender := services.SMTPSender{}

	go services.StartCronJob(db, emailSender)

	r.Run(":8080")
}
