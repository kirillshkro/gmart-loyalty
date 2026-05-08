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
	Create(order model.Order) error
}

type Getter interface {
	GetByID(ctx context.Context, id int) (model.Order, error)
	GetAll(ctx context.Context) ([]model.Order, error)
}

func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (o OrderRepository) Create(order model.Order) error {
	th := o.onConflict()
	err := th.Transaction(func(tx *gorm.DB) error {
		if err := gorm.G[model.Order](tx).Create(context.Background(), &order); err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				anotherUser, err := gorm.G[model.UserProfile](tx).Select("user_id").Where("order_num=?", order.OrderNum).First(context.Background())
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

func (o OrderRepository) GetByID(ctx context.Context, id int) (model.Order, error) {
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

func (o OrderRepository) GetAll(ctx context.Context) ([]model.Order, error) {
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

func (o OrderRepository) onConflict() *gorm.DB {
	return o.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "order_num"}},
		DoNothing: true,
	},
		clause.Returning{Columns: []clause.Column{{Name: "id"}}})
}
