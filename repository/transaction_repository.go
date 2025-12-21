package repository

import (
	"github.com/IndalAwalaikal/warung-pos/backend/config"
	"github.com/IndalAwalaikal/warung-pos/backend/model"
	"gorm.io/gorm"
)

type TransactionRepository interface {
    Create(tx *model.Transaction) error
    List() ([]model.Transaction, error)
    GetByID(id uint) (*model.Transaction, error)
}

type transactionRepo struct{
    db *gorm.DB
}

func NewTransactionRepository() TransactionRepository {
    return &transactionRepo{db: config.DB}
}

func (r *transactionRepo) Create(tx *model.Transaction) error {
    return r.db.Create(tx).Error
}

func (r *transactionRepo) List() ([]model.Transaction, error) {
    var list []model.Transaction
    if err := r.db.Preload("Items.Menu").Find(&list).Error; err != nil {
        return nil, err
    }
    return list, nil
}

func (r *transactionRepo) GetByID(id uint) (*model.Transaction, error) {
    var t model.Transaction
    if err := r.db.Preload("Items.Menu").First(&t, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, err
    }
    return &t, nil
}
