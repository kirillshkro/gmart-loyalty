package model

import "time"

type Order struct {
	ID         int    `gorm:"primaryKey"`
	Number     string `json:"number" gorm:"index;unique;type:varchar(50);not null"`
	UserID     int
	Accrual    int         `json:"accrual"`
	UploadedAt time.Time   `gorm:"autoCreateTime"`
	Balance    UserBalance `gorm:"foreignKey:OrderID;references:ID;constraints:OnDelete:SET NULL"`
	User       UserProfile `json:"-"`
}
