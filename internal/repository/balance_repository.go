package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/kirillshkro/gmart-loyalty/internal/model"
	"github.com/kirillshkro/gmart-loyalty/internal/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IBalanceRepository interface {
	BalanceByUser(userID int) (model.UserBalance, error)
	SetWithdraw(ctx context.Context, withdraw *types.WithdrawRequest) error
	ListWithdraws(userID int) ([]types.WithdrawResponse, error)
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
	userID, ok := ctx.Value(types.UserID).(int)
	if !ok {
		return fmt.Errorf("user_id not found in context or invalid type")
	}
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

func (b *Repository) ListWithdraws(userID int) ([]types.WithdrawResponse, error) {
	var (
		withdraws []types.WithdrawResponse
		withdraw  types.WithdrawResponse
	)
	balances, err := gorm.G[model.UserBalance](b.db).Select("order_number", "withdrawn", "updated_at").Where("user_id=?", userID).Order(
		clause.OrderByColumn{
			Desc:   true,
			Column: clause.Column{Name: "updated_at"},
		},
	).Find(context.Background())
	if err != nil {
		return nil, err
	}
	for _, balance := range balances {
		withdraw.Order = balance.OrderNumber
		withdraw.ProcessedAt = balance.UpdatedAt
		withdraw.Sum = int(balance.Withdrawn)
		withdraws = append(withdraws, withdraw)
	}
	return withdraws, nil
}
