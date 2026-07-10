package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName   string `json:"first_name" gorm:"size:100"`
	LastName    string `json:"last_name" gorm:"size:100"`
	PhoneNumber string `json:"phone_number" gorm:"unique;size:20;index"`
	Email       string `json:"email" gorm:"unique;not null;size:255;index"`
	Password    string `json:"-" gorm:"not null;size:255"`
}
