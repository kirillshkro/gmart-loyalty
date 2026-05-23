package repository

import (
	"context"
	"errors"

	"github.com/kirillshkro/gmart-loyalty/internal/model"
	"github.com/kirillshkro/gmart-loyalty/internal/types"
	"gorm.io/gorm"
)

type IBalanceRepository interface {
	BalanceByUser(userID int) (model.UserBalance, error)
	SetWithdraw(ctx context.Context, withdraw *types.WithdrawRequest) error
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

func (b Repository) SetWithdraw(ctx context.Context, withdraw *types.WithdrawRequest) error {
	userID := ctx.Value(types.UserID).(int)
	err := b.db.Transaction(func(tx *gorm.DB) error {
		balance, err := gorm.G[model.UserBalance](b.db).Where("user_id =? and order_number =?", userID, withdraw.Order).First(ctx)
		if err != nil {
			return err
		}
		if balance.Current < float64(withdraw.Sum) {
			return &types.ErrInsufficientBalance{
				UserID:  userID,
				Balance: balance.Current,
			}
		}
		balance.Current -= float64(withdraw.Sum)
		balance.Withdrawn = float64(withdraw.Sum)
		rows, err := gorm.G[model.UserBalance](b.db).Where("user_id = ? and order_number =?", userID, withdraw.Order).Updates(ctx, balance)
		if err != nil {
			return err
		}
		if rows == 0 {
			return errors.New("balance not updated")
		}
		return nil
	})
	return err
}
