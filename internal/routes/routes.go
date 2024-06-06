package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"rutube-assignment/internal/controllers"
	"rutube-assignment/internal/middlewares"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	userController := &controllers.UserController{DB: db}
	authController := &controllers.AuthController{DB: db}

	r.POST("/register", authController.Register)
	r.POST("/login", authController.Login)

	protected := r.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/users", userController.GetAllUsers)
		protected.POST("/subscribe", userController.Subscribe)
		protected.POST("/unsubscribe", userController.Unsubscribe)
	}
}
