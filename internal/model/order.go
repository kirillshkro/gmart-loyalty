// internal/model/order.go
package model

import (
	"time"

	"gorm.io/gorm"
)

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
	Accrual    float64     `json:"accrual" gorm:"type:decimal(10,2); not null"`
	Status     OrderStatus `json:"status" gorm:"type:enum;not null;index"`
	UploadedAt time.Time   `gorm:"autoCreateTime"`
	Balance    UserBalance `gorm:"foreignKey:OrderID;references:ID;constraints:OnDelete:SET NULL"`
	User       UserProfile `json:"-"`
}

// Добавляет в таблицу UserBalance запись с текущим балансом и отсутствием выведенных средств
// К текущему балансу (Current) добавляется величина в поле Accrual
func (o Order) AfterCreate(tx *gorm.DB) error {
	var (
		balance        UserBalance
		currentBalance float64
	)

	// Получаем предыдущий баланс пользователя
	if err := tx.Model(&UserBalance{}).Where("user_id = ?", o.UserID).Select("COALESCE(SUM(current), 0)").Scan(&currentBalance).Error; err != nil {
		return err
	}

	// Создаем новую запись в UserBalance
	balance.OrderID = o.ID
	balance.UserID = o.UserID
	balance.Current = currentBalance + o.Accrual
	balance.Withdrawn = 0

	if err := tx.Create(&balance).Error; err != nil {
		return err
	}
	return nil
}
