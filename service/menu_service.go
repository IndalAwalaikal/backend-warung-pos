package service

import (
	"github.com/IndalAwalaikal/warung-pos/backend/model"
	"github.com/IndalAwalaikal/warung-pos/backend/repository"
)

type MenuService interface {
    Create(m *model.Menu) error
    List() ([]model.Menu, error)
    GetByID(id uint) (*model.Menu, error)
    Update(m *model.Menu) error
    Delete(id uint) error
}

type menuService struct{
    repo repository.MenuRepository
}

func NewMenuService(r repository.MenuRepository) MenuService {
    return &menuService{repo: r}
}

func (s *menuService) Create(m *model.Menu) error {
    return s.repo.Create(m)
}

func (s *menuService) List() ([]model.Menu, error) {
    return s.repo.List()
}

func (s *menuService) GetByID(id uint) (*model.Menu, error) {
    return s.repo.GetByID(id)
}

func (s *menuService) Update(m *model.Menu) error {
    return s.repo.Update(m)
}

func (s *menuService) Delete(id uint) error {
    return s.repo.Delete(id)
}
