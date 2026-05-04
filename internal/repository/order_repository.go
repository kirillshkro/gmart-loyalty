package repository

import (
	"context"

	"github.com/kirillshkro/gmart-loyalty/internal/model"
	"gorm.io/gorm"
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
