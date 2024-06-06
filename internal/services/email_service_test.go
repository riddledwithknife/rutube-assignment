package services

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"rutube-assignment/internal/config"
	"rutube-assignment/internal/models"
	"testing"
	"time"
)

type MockEmailSender struct {
	SentEmails []MockEmail
}

type MockEmail struct {
	To      string
	Subject string
	Body    string
}

func (m *MockEmailSender) Send(to string, subject string, body string) error {
	m.SentEmails = append(m.SentEmails, MockEmail{To: to, Subject: subject, Body: body})
	return nil
}

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

func TestCheckAndSendBirthdayGreetings(t *testing.T) {
	mockEmailSender := &MockEmailSender{}

	birthday := time.Now().Format("2006-01-02")
	user := models.User{
		Name:     "Birthday User",
		Email:    "birthday@example.com",
		Password: "password123",
		Birthday: birthday,
	}
	db.Create(&user)

	subscriber := models.User{
		Name:     "Subscriber User",
		Email:    "subscriber@example.com",
		Password: "password123",
		Birthday: "2000-01-01",
	}
	db.Create(&subscriber)

	db.Create(&models.Subscription{
		UserID:       subscriber.ID,
		SubscribedTo: user.ID,
	})

	CheckAndSendBirthdayGreetings(db, mockEmailSender)

	assert.Len(t, mockEmailSender.SentEmails, 1)
	sentEmail := mockEmailSender.SentEmails[0]
	assert.Equal(t, "subscriber@example.com", sentEmail.To)
	assert.Equal(t, "Birthday Notification", sentEmail.Subject)
	assert.Equal(t, "Today is Birthday User's birthday!", sentEmail.Body)
}
