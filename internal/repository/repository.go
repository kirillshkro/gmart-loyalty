package repository

import "gorm.io/gorm"

type IRepository interface {
	IUserRepository
	IOrderRepository
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) IRepository {
	return &Repository{
		db: db,
	}
}
