package model

type UserBalance struct {
	ID        int     `gorm:"primaryKey"`
	UserID    int     `gorm:"index;not null"`
	OrderID   int     `gorm:"index;not null"`
	Current   float64 `json:"current" gorm:"not null;type:decimal(10,2)"`
	Withdrawn float64 `json:"withdrawn" gorm:"type:decimal(10,2)"`
}
