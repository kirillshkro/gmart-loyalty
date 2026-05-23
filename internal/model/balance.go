package model

import "time"

type UserBalance struct {
	ID          int       `gorm:"primaryKey"`
	UserID      int       `gorm:"index;not null;index:user_id_order"`
	OrderNumber string    `gorm:"index;not null;index:user_id_order"`
	Current     float64   `json:"current" gorm:"not null;type:decimal(10,2)"`
	Withdrawn   float64   `json:"withdrawn" gorm:"type:decimal(10,2)"`
	UpdatedAt   time.Time `gorm:"not null; index; autoUpdateTime"`
}
