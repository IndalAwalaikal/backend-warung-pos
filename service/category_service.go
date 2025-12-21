package service

import (
	"github.com/IndalAwalaikal/warung-pos/backend/model"
	"github.com/IndalAwalaikal/warung-pos/backend/repository"
)

type CategoryService interface {
    Create(cat *model.Category) error
    List() ([]model.Category, error)
}

type categoryService struct{
    repo repository.CategoryRepository
}

func NewCategoryService(r repository.CategoryRepository) CategoryService {
    return &categoryService{repo: r}
}

func (s *categoryService) Create(cat *model.Category) error {
    return s.repo.Create(cat)
}

func (s *categoryService) List() ([]model.Category, error) {
    return s.repo.List()
}
