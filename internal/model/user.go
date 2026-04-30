package model

type User struct {
	Name     string `json:"name" gorm:"index;not null;unique"`
	Password string `json:"password" gorm:"not null"`
}

type UserProfile struct {
	ID   int `gorm:"primaryKey"`
	User `gorm:"embedded"`
}
