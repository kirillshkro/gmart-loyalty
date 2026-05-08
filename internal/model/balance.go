package model

type UserBalance struct {
	ID      int `gorm:"primaryKey"`
	UserID  int `gorm:"index;not null"`
	OrderID int `gorm:"index;not null"`
	Balance int `gorm:"not null"`
}
