package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"rutube-assignment/internal/config"
	"rutube-assignment/internal/middlewares"
	"rutube-assignment/internal/models"
	"rutube-assignment/internal/services"
	"rutube-assignment/internal/testutils"
	"strings"
	"testing"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	var err error
	db, err = gorm.Open(postgres.Open(config.GetDBConnectionString()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database")
	}

	db.AutoMigrate(&models.User{}, &models.Subscription{})

	code := m.Run()

	os.Exit(code)
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	return r
}

func TestSubscribe(t *testing.T) {
	r := SetupRouter()
	userController := UserController{DB: db}
	r.POST("/subscribe", middlewares.AuthMiddleware(), userController.Subscribe)

	db.Create(&models.User{
		Name:     "Maxim User",
		Email:    "maxim@example.ru",
		Password: "password123",
		Birthday: "2000-01-01",
	})

	var hashedPassword, _ = services.HashPassword("password123")

	testUser := models.User{
		Name:     "Andrey User",
		Email:    "andrey@example.ru",
		Password: hashedPassword,
		Birthday: "2000-01-01",
	}
	db.Create(&testUser)

	token, err := testutils.GenerateMockToken(testUser.ID)
	if err != nil {
		t.Fatalf("Failed to generate mock token: %v", err)
	}

	subscribeRequest := `{
        "email": "maxim@example.com"
    }`

	req, _ := http.NewRequest("POST", "/subscribe", strings.NewReader(subscribeRequest))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Successfully subscribed to birthday notifications")
}

func TestUnsubscribe(t *testing.T) {
	r := SetupRouter()
	userController := UserController{DB: db}
	r.POST("/unsubscribe", middlewares.AuthMiddleware(), userController.Unsubscribe)

	jane := models.User{
		Name:     "Anton User",
		Email:    "anton@example.com",
		Password: "password123",
		Birthday: "2000-01-01",
	}
	db.Create(&jane)

	var hashedPassword, _ = services.HashPassword("password123")

	testUser := models.User{
		Name:     "Oleg User",
		Email:    "oleg@example.com",
		Password: hashedPassword,
		Birthday: "2000-01-01",
	}
	db.Create(&testUser)

	db.Create(&models.Subscription{
		UserID:       testUser.ID,
		SubscribedTo: jane.ID,
	})

	token, err := testutils.GenerateMockToken(testUser.ID)
	if err != nil {
		t.Fatalf("Failed to generate mock token: %v", err)
	}

	unsubscribeRequest := `{
        "email": "jane@example.com"
    }`

	req, _ := http.NewRequest("POST", "/unsubscribe", strings.NewReader(unsubscribeRequest))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Successfully unsubscribed from birthday notifications")
}
