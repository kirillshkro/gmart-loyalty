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

type OrderRepository struct {
	db *gorm.DB
}

type IOrderRepository interface {
	Setter
	Getter
}

type Setter interface {
	Create(ctx context.Context, order model.Order) error
}

type Getter interface {
	GetByID(id int) (model.Order, error)
	GetAll() ([]model.Order, error)
}

func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (o OrderRepository) Create(ctx context.Context, order model.Order) error {
	th := o.onConflict()
	err := th.Transaction(func(tx *gorm.DB) error {
		if err := gorm.G[model.Order](tx).Create(ctx, &order); err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				anotherUser, err := gorm.G[model.UserProfile](tx).Select("user_id").Where("order_num=?", order.OrderNum).First(ctx)
				if err != nil {
					return fmt.Errorf("user %d can't own order %s", anotherUser.ID, order.OrderNum)
				}
				return &types.ErrOwnAnotherUser{
					UserID:   anotherUser.ID,
					OrderNum: order.OrderNum,
				}
			}
			return err
		}
		return nil
	})
	return err
}

func (o OrderRepository) GetByID(id int) (model.Order, error) {
	var (
		order model.Order
		err   error
	)
	if order, err = gorm.G[model.Order](o.db).Where("id = ?", id).First(context.Background()); err != nil {
		return model.Order{}, err
	}
	return order, nil
}

func (o OrderRepository) GetAll() ([]model.Order, error) {
	var (
		orders []model.Order
		err    error
	)
	if orders, err = gorm.G[model.Order](o.db).Order(
		clause.OrderByColumn{
			Desc:   false,
			Column: clause.Column{Name: "created_at"},
		},
	).Find(context.Background()); err != nil {
		return nil, err
	}
	return orders, nil

}

func (o OrderRepository) onConflict() *gorm.DB {
	return o.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "order_num"}},
		DoNothing: true,
	},
		clause.Returning{Columns: []clause.Column{{Name: "id"}}})
}
