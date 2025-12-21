package repository

import (
	"github.com/IndalAwalaikal/warung-pos/backend/config"
	"github.com/IndalAwalaikal/warung-pos/backend/model"
	"gorm.io/gorm"
)

type CategoryRepository interface {
    Create(cat *model.Category) error
    List() ([]model.Category, error)
}

type categoryRepo struct{
    db *gorm.DB
}

func NewCategoryRepository() CategoryRepository {
    return &categoryRepo{db: config.DB}
}

func (r *categoryRepo) Create(cat *model.Category) error {
    return r.db.Create(cat).Error
}

func (r *categoryRepo) List() ([]model.Category, error) {
    var list []model.Category
    if err := r.db.Find(&list).Error; err != nil {
        return nil, err
    }
    return list, nil
}
