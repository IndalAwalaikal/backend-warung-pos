package service

import (
	"github.com/IndalAwalaikal/warung-pos/backend/model"
	"github.com/IndalAwalaikal/warung-pos/backend/repository"
)

type TransactionService interface {
    Create(tx *model.Transaction) error
    List() ([]model.Transaction, error)
    GetByID(id uint) (*model.Transaction, error)
}

type transactionService struct{
    repo repository.TransactionRepository
}

func NewTransactionService(r repository.TransactionRepository) TransactionService {
    return &transactionService{repo: r}
}

func (s *transactionService) Create(tx *model.Transaction) error {
    return s.repo.Create(tx)
}

func (s *transactionService) List() ([]model.Transaction, error) {
    return s.repo.List()
}

func (s *transactionService) GetByID(id uint) (*model.Transaction, error) {
    return s.repo.GetByID(id)
}
