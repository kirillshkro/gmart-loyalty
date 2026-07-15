package repository

import (
	"context"
	"fmt"

	"github.com/kirillshkro/gmart-loyalty/internal/model"
	"github.com/kirillshkro/gmart-loyalty/internal/types"
	"gorm.io/gorm"
)

type IBalanceRepository interface {
	BalanceByUser(userID int) (model.UserBalance, error)
	SetWithdraw(ctx context.Context, withdraw *types.WithdrawRequest) error
	ListWithdraws(userID int) ([]types.WithdrawResponse, error)
}

func (b Repository) BalanceByUser(userID int) (model.UserBalance, error) {
	var balance model.UserBalance
	err := b.db.Where("user_id = ?", userID).First(&balance).Error
	if err != nil {
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
		// Получаем общий баланс пользователя
		var userBalance model.UserBalance
		err := tx.Where("user_id = ?", userID).First(&userBalance).Error
		if err != nil {
			return err
		}

		// Проверяем достаточно ли средств на общем балансе
		if userBalance.Current < float64(withdraw.Sum) {
			return &types.ErrInsufficientBalance{
				UserID:  userID,
				Balance: userBalance.Current,
			}
		}

		// Создаем запись о выводе средств
		withdrawal := model.Withdrawal{
			UserID:      userID,
			OrderNumber: withdraw.Order,
			Sum:         float64(withdraw.Sum),
		}

		if err := tx.Create(&withdrawal).Error; err != nil {
			return err
		}

		// Обновляем баланс пользователя
		userBalance.Current -= float64(withdraw.Sum)
		userBalance.Withdrawn += float64(withdraw.Sum)

		if err := tx.Save(&userBalance).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func (b Repository) ListWithdraws(userID int) ([]types.WithdrawResponse, error) {
	var withdrawals []types.WithdrawResponse

	// Получаем все записи о выводах для пользователя
	var dbWithdrawals []model.Withdrawal
	err := b.db.Where("user_id = ?", userID).Order("processed_at DESC").Find(&dbWithdrawals).Error
	if err != nil {
		return nil, err
	}

	// Преобразуем в ответные структуры
	for _, dbWithdrawal := range dbWithdrawals {
		withdrawal := types.WithdrawResponse{
			Order:       dbWithdrawal.OrderNumber,
			ProcessedAt: dbWithdrawal.ProcessedAt,
			Sum:         int(dbWithdrawal.Sum),
		}
		withdrawals = append(withdrawals, withdrawal)
	}

	return withdrawals, nil
}
