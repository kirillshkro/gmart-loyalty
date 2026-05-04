package model

type UserBalance struct {
	ID      int     `gorm:"primaryKey"`
	UserID  int     `gorm:"index;not null"`
	OrderID int     `gorm:"index;not null"`
	Balance float64 `gorm:"not null;type:decimal(10,2)"`
}
