package model

import "time"

type Order struct {
	ID        int `gorm:"primaryKey"`
	OrderID   int `gorm:"index;unique"`
	UserID    int
	CreatedAt time.Time
	Balance   UserBalance `gorm:"foreignKey:OrderID;references:ID;constraints:OnDelete:SET NULL"`
}
