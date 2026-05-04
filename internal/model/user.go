package model

type User struct {
	UserName  string `json:"user_name" gorm:"index;not null;unique"`
	Password  string `json:"password" gorm:"not null"`
	Password2 string `json:"password2" gorm:"-:all"`
}

type UserProfile struct {
	ID     int `gorm:"primaryKey"`
	User   `gorm:"embedded"`
	Orders []Order `gorm:"foreignKey:UserID;references:ID;constraints:OnDelete:CASCADE"`
}
