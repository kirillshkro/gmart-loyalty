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
	if err := gorm.G[model.Order](th).Create(ctx, order); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			ok := o.anotherUser(ctx)
			if ok {
				return &types.ErrOwnAnotherUser{
					UserID:   order.UserID,
					OrderNum: order.OrderNum,
				}
			}
		}
		return err
	}
	return nil
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
	userID := ctx.Value("user_id").(int)
	if orders, err = gorm.G[model.Order](o.db).Order(
		clause.OrderByColumn{
			Desc:   false,
			Column: clause.Column{Name: "created_at"},
		},
	).Where("user_id = ?", userID).Find(ctx); err != nil {
		return nil, err
	}
	return orders, nil

}

func (o Repository) onOrderConflict() *gorm.DB {
	return o.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "order_num"}},
		DoNothing: true,
	},
		clause.Returning{Columns: []clause.Column{{Name: "id"}}})
}

func (o Repository) anotherUser(ctx context.Context) bool {
	//извлечь user_id
	userID := ctx.Value(types.UserID).(int)
	//извлечь номер заказа
	numOrder := ctx.Value(types.OrderNum).(string)

	order, _ := gorm.G[model.Order](o.db).Joins(clause.JoinTarget{
		Table: "orders",
		Type:  clause.InnerJoin,
	}, nil).Select("user_id").Where("order_num = ?", numOrder).First(ctx)

	return userID != order.UserID
}
