package repository

import (
	"context"

	"github.com/kirillshkro/gmart-loyalty/internal/model"
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
	GetByID(id int) (model.Order, error)
	GetAll() ([]model.Order, error)
}

func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (o OrderRepository) Create(order model.Order) error {
	if err := gorm.G[model.Order](o.db).Create(context.Background(), &order); err != nil {
		return err
	}
	return nil
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
