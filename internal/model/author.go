package model

type Author struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	FullName string `json:"full_name"`
	Bio      string `json:"bio"`
}
