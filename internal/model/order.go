package model

import "time"

type OrderStatus string

const (
	StatusNew        = "New"
	StatusProcessing = "Processing"
	StutusInvalid    = "Invalid"
	StatusProcessed  = "Processed"
)

type Order struct {
	ID         int    `gorm:"primaryKey"`
	Number     string `json:"number" gorm:"index;unique;type:varchar(50);not null"`
	UserID     int
	Accrual    int         `json:"accrual"`
	Status     OrderStatus `json:"status" gorm:"type:enum;not null;index"`
	UploadedAt time.Time   `gorm:"autoCreateTime"`
	Balance    UserBalance `gorm:"foreignKey:OrderID;references:ID;constraints:OnDelete:SET NULL"`
	User       UserProfile `json:"-"`
}
