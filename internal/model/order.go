package model

import "time"

type Order struct {
	ID        int    `gorm:"primaryKey"`
	OrderNum  string `gorm:"index;unique;type:varchar(50);not null"`
	UserID    int
	CreatedAt time.Time
	Balance   UserBalance `gorm:"foreignKey:OrderID;references:ID;constraints:OnDelete:SET NULL"`
}
