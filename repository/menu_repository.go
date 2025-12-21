package repository

import (
	"github.com/IndalAwalaikal/warung-pos/backend/config"
	"github.com/IndalAwalaikal/warung-pos/backend/model"
	"gorm.io/gorm"
)

type MenuRepository interface {
    Create(m *model.Menu) error
    List() ([]model.Menu, error)
    GetByID(id uint) (*model.Menu, error)
    Update(m *model.Menu) error
    Delete(id uint) error
}

type menuRepo struct{
    db *gorm.DB
}

func NewMenuRepository() MenuRepository {
    return &menuRepo{db: config.DB}
}

func (r *menuRepo) Create(m *model.Menu) error {
    return r.db.Create(m).Error
}

func (r *menuRepo) List() ([]model.Menu, error) {
    var list []model.Menu
    if err := r.db.Preload("Category").Find(&list).Error; err != nil {
        return nil, err
    }
    return list, nil
}

func (r *menuRepo) GetByID(id uint) (*model.Menu, error) {
    var m model.Menu
    if err := r.db.Preload("Category").First(&m, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, err
    }
    return &m, nil
}

func (r *menuRepo) Update(m *model.Menu) error {
    return r.db.Save(m).Error
}

func (r *menuRepo) Delete(id uint) error {
    return r.db.Delete(&model.Menu{}, id).Error
}
