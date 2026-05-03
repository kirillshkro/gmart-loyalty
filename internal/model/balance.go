package model

type UserBalance struct {
	ID      int     `gorm:"primaryKey"`
	UserID  int     `gorm:"index;not null"`
	Balance float64 `gorm:"type:decimal(10,2)"`
}
