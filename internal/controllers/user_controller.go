package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"rutube-assignment/internal/models"
)

type UserController struct {
	DB *gorm.DB
}

func (uc *UserController) GetAllUsers(c *gin.Context) {
	var users []models.User
	uc.DB.Find(&users)
	c.JSON(http.StatusOK, users)
}

func (uc *UserController) Subscribe(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var subscribedToUser models.User
	if err := uc.DB.Where("email = ?", request.Email).First(&subscribedToUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve user ID from context"})
		return
	}

	var existingSubscription models.Subscription
	if err := uc.DB.Where("user_id = ? AND subscribed_to = ?", userID, subscribedToUser.ID).First(&existingSubscription).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already subscribed"})
		return
	}

	subscription := models.Subscription{
		UserID:       userID.(uint),
		SubscribedTo: subscribedToUser.ID,
	}
	if err := uc.DB.Create(&subscription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully subscribed to birthday notifications for " + subscribedToUser.Name})
}

func (uc *UserController) Unsubscribe(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var subscribedToUser models.User
	if err := uc.DB.Where("email = ?", request.Email).First(&subscribedToUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve user ID from context"})
		return
	}

	if err := uc.DB.Where("user_id = ? AND subscribed_to = ?", userID, subscribedToUser.ID).Delete(&models.Subscription{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully unsubscribed from birthday notifications for " + subscribedToUser.Name})
}
