package model

type Order struct {
	ID      int `gorm:"primaryKey"`
	OrderID int `gorm:"index;unique"`
	UserID  int
}
