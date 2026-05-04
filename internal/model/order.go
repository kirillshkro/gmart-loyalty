package model

type Order struct {
	ID      int `gorm:"primaryKey"`
	OrderID int `gorm:"index;unique"`
	UserID  int
	Balance UserBalance `gorm:"foreignKey:OrderID;references:ID;constraints:OnDelete:SET NULL"`
}
