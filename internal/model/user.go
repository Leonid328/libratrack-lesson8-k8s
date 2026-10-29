package model

type User struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	PasswordHash string `json:"-"`
}
