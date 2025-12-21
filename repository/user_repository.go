package repository

import (
	"errors"

	"github.com/IndalAwalaikal/warung-pos/backend/config"
	"github.com/IndalAwalaikal/warung-pos/backend/model"
	"gorm.io/gorm"
)

type UserRepository interface {
    Create(user *model.User) error
    FindByEmail(email string) (*model.User, error)
    FindByID(id uint) (*model.User, error)
}

type userRepo struct{
    db *gorm.DB
}

func NewUserRepository() UserRepository {
    return &userRepo{db: config.DB}
}

func (r *userRepo) Create(user *model.User) error {
    return r.db.Create(user).Error
}

func (r *userRepo) FindByEmail(email string) (*model.User, error) {
    var u model.User
    if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &u, nil
}

func (r *userRepo) FindByID(id uint) (*model.User, error) {
    var u model.User
    if err := r.db.First(&u, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &u, nil
}
