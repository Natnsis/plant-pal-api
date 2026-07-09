package models

type User struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number" gorm:"unique"`
	Email       string `json:"email" gorm:"unique;not null"`
	Password    string `json:"password" gorm:"not null"`
}
