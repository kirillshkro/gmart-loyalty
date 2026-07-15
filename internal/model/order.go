package model

import (
	"time"

	"gorm.io/gorm"
)

type OrderStatus string

const (
	StatusNew        OrderStatus = "NEW"
	StatusProcessing OrderStatus = "PROCESSING"
	StatusInvalid    OrderStatus = "INVALID"
	StatusProcessed  OrderStatus = "PROCESSED"
)

type Order struct {
	ID         int    `gorm:"primaryKey"`
	Number     string `json:"number" gorm:"index;unique;type:varchar(50);not null"`
	UserID     int
	Accrual    float64     `json:"accrual" gorm:"type:decimal(10,2); not null"`
	Status     OrderStatus `json:"status" gorm:"not null;index"`
	UploadedAt time.Time   `gorm:"autoCreateTime;index"`
	User       UserProfile `gorm:"foreignKey:UserID"`
}

// AfterCreate добавляет запись в баланс пользователя после создания заказа
func (o Order) AfterCreate(tx *gorm.DB) error {
	var (
		balance        UserBalance
		currentBalance float64
	)

	// Получаем текущий баланс пользователя
	if err := tx.Model(&UserBalance{}).Where("user_id = ?", o.UserID).Select("COALESCE(SUM(current), 0)").Scan(&currentBalance).Error; err != nil {
		return err
	}

	// Создаем новую запись в UserBalance
	balance.UserID = o.UserID
	balance.OrderNumber = o.Number
	balance.Current = currentBalance + o.Accrual
	balance.Withdrawn = 0

	if err := tx.Create(&balance).Error; err != nil {
		return err
	}
	return nil
}
