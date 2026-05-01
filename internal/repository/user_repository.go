package repository

import (
	"context"
	"errors"

	"github.com/kirillshkro/gmart-loyalty/internal/model"
	"github.com/kirillshkro/gmart-loyalty/internal/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	th := u.onConflict()
	if err := gorm.G[model.UserProfile](th).Create(context.Background(), &userProfile); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return &types.ErrDuplicateUser{
				UserName: userProfile.UserName,
			}
		}
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

func (u UserRepository) onConflict() *gorm.DB {
	return u.db.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_name"}},
			DoNothing: true,
		},
		clause.Returning{Columns: []clause.Column{{Name: "id"}}},
	)
}
