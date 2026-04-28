package repository

import (
	"context"

	"github.com/kirillshkro/gmart-loyalty/internal/model"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Create(profile model.UserProfile) error
	GetByID(id int) (model.UserProfile, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{
		db: db,
	}
}
func (u *UserRepository) Create(userProfile model.UserProfile) error {
	if err := gorm.G[model.UserProfile](u.db).Create(context.Background(), &userProfile); err != nil {
		return err
	}
	return nil
}

func (u UserRepository) GetByID(id int) (model.UserProfile, error) {
	var (
		up  model.UserProfile
		err error
	)

	if up, err = gorm.G[model.UserProfile](u.db).Where("id = ?", id).First(context.Background()); err != nil {
		return model.UserProfile{}, err
	}
	return up, nil
}
