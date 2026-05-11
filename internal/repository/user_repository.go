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
	CreateUser(profile model.UserProfile) (int, error)
	UserByID(id int) (model.UserProfile, error)
	UserByName(username string) (model.UserProfile, error)
}

type UserRepository struct {
	db *gorm.DB
}

func (u *Repository) CreateUser(userProfile model.UserProfile) (int, error) {
	th := u.onConflict()
	if err := gorm.G[model.UserProfile](th).Create(context.Background(), &userProfile); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return 0, &types.ErrDuplicateUser{
				UserName: userProfile.UserName,
			}
		}
		return 0, err
	}
	return userProfile.ID, nil
}

func (u Repository) UserByID(id int) (model.UserProfile, error) {
	var (
		up  model.UserProfile
		err error
	)

	if up, err = gorm.G[model.UserProfile](u.db).Where("id = ?", id).First(context.Background()); err != nil {
		return model.UserProfile{}, err
	}
	return up, nil
}

func (u Repository) UserByName(userName string) (model.UserProfile, error) {
	var (
		up  model.UserProfile
		err error
	)
	if up, err = gorm.G[model.UserProfile](u.db).Where("user_name = ?", userName).First(context.Background()); err != nil {
		return model.UserProfile{}, err
	}
	return up, nil
}

func (u Repository) onConflict() *gorm.DB {
	return u.db.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_name"}},
			DoNothing: true,
		},
		clause.Returning{Columns: []clause.Column{{Name: "id"}}},
	)
}
