package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FullName    string `json:"full_name" gorm:"size:100"`
	Email       string `json:"email" gorm:"unique;not null;size:255;index"`
	PhoneNumber string `json:"phone_number" gorm:"unique;not null"`
	Password    string `json:"-" gorm:"not null;size:255"`
}
