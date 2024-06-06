package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	Birthday string `json:"birthday"`
}

type Subscription struct {
	gorm.Model
	UserID       uint `json:"user_id"`
	SubscribedTo uint `json:"subscribed_to"`
}
