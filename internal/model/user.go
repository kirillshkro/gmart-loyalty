package model

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserProfile struct {
	ID   int `gorm:"primaryKey"`
	User `gorm:"embedded"`
}
