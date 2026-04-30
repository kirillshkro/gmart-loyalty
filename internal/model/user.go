package model

type User struct {
	Name      string `json:"name" gorm:"index;not null;unique"`
	Password  string `json:"password" gorm:"not null"`
	Password2 string `json:"password2" gorm:"-:all"`
}

type UserProfile struct {
	ID   int `gorm:"primaryKey"`
	User `gorm:"embedded"`
}
