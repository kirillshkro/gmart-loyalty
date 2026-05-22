package repository

import (
	"context"

	"github.com/kirillshkro/gmart-loyalty/internal/model"
	"gorm.io/gorm"
)

type IBalanceRepository interface {
	BalanceByUser(userID int) (model.UserBalance, error)
}

func (b Repository) BalanceByUser(userID int) (model.UserBalance, error) {
	var (
		balance model.UserBalance
		err     error
	)
	if balance, err = gorm.G[model.UserBalance](b.db).Where("user_id = ?", userID).First(context.Background()); err != nil {
		return balance, err
	}
	return balance, nil
}
