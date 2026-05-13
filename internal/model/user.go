package model

type User struct {
	Login     string `json:"login" gorm:"index;not null;unique"`
	Password  string `json:"password" gorm:"not null"`
	Password2 string `json:"password2" gorm:"-:all"`
}

type UserProfile struct {
	ID     int `gorm:"primaryKey"`
	User   `gorm:"embedded"`
	Orders []Order `gorm:"foreignKey:UserID;references:ID;constraints:OnDelete:CASCADE"`
}
