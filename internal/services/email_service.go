package services

import (
	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"log"
	"rutube-assignment/internal/models"
	"time"
)

type EmailSender interface {
	Send(to string, subject string, body string) error
}

type SMTPSender struct{}

func (s SMTPSender) Send(to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "youremail@example.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.example.com", 587, "user", "password")

	return d.DialAndSend(m)
}

func CheckAndSendBirthdayGreetings(db *gorm.DB, emailSender EmailSender) {
	var users []models.User
	today := time.Now().Format("2006-01-02")
	db.Where("birthday = ?", today).Find(&users)

	for _, user := range users {
		var subscriptions []models.Subscription
		db.Where("subscribed_to = ?", user.ID).Find(&subscriptions)

		for _, subscription := range subscriptions {
			var subscriber models.User
			db.First(&subscriber, subscription.UserID)
			err := emailSender.Send(subscriber.Email, "Birthday Notification", "Today is "+user.Name+"'s birthday!")
			if err != nil {
				log.Println("Could not send email:", err)
			}
		}
	}
}

func StartCronJob(db *gorm.DB, emailSender EmailSender) {
	c := cron.New()
	c.AddFunc("@daily", func() {
		CheckAndSendBirthdayGreetings(db, emailSender)
	})
	c.Start()
}
