package repository

import (
	"context"
	"errors"

	"github.com/kirillshkro/gmart-loyalty/internal/model"
	"github.com/kirillshkro/gmart-loyalty/internal/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IOrderRepository interface {
	Setter
	Getter
}

type Setter interface {
	CreateOrder(ctx context.Context, order *model.Order) error
}

type Getter interface {
	OrderByID(ctx context.Context, id int) (model.Order, error)
	GetAll(ctx context.Context) ([]model.Order, error)
}

func (o Repository) CreateOrder(ctx context.Context, order *model.Order) error {
	th := o.onOrderConflict()
	err := o.db.Transaction(func(tx *gorm.DB) error {
		if err := gorm.G[model.Order](th).Create(ctx, order); err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				ok := o.anotherUser(ctx)
				if ok {
					return &types.ErrOwnAnotherUser{
						UserID:   order.UserID,
						OrderNum: order.Number,
					}
				}
			}
			return err
		}
		order.Status = model.StatusProcessed
		rows, err := gorm.G[model.Order](tx).Where("id = ? AND user_id = ?", order.ID, order.UserID).Update(ctx, "status", order.Status)
		if err != nil {
			return err
		}
		if rows == 0 {
			return err
		}
		return nil
	})
	return err
}

func (o Repository) OrderByID(ctx context.Context, id int) (model.Order, error) {
	var (
		order model.Order
		err   error
	)
	userID := ctx.Value("user_id").(int)
	if order, err = gorm.G[model.Order](o.db).Where("id = ? AND user_id = ?", id, userID).First(ctx); err != nil {
		return model.Order{}, err
	}
	return order, nil
}

func (o Repository) GetAll(ctx context.Context) ([]model.Order, error) {
	var (
		orders []model.Order
		err    error
	)
	userID := ctx.Value(types.UserID).(int)
	if orders, err = gorm.G[model.Order](o.db).Order(
		clause.OrderByColumn{
			Desc:   true,
			Column: clause.Column{Name: "uploaded_at"},
		},
	).Where("user_id = ?", userID).Find(ctx); err != nil {
		return nil, err
	}
	return orders, nil
}

func (o Repository) onOrderConflict() *gorm.DB {
	return o.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "number"}},
		DoNothing: true,
	},
		clause.Returning{Columns: []clause.Column{{Name: "id"}}})
}

func (o Repository) anotherUser(ctx context.Context) bool {
	//извлечь user_id
	userID := ctx.Value(types.UserID).(int)
	//извлечь номер заказа
	numOrder := ctx.Value(types.OrderNum).(string)

	order, _ := gorm.G[model.Order](o.db).Select("user_id").Where("number = ?", numOrder).First(ctx)

	return userID != order.UserID
}
