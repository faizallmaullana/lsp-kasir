package services

import (
	"errors"
	"faizalmaulana/lsp/models/entity"
	"faizalmaulana/lsp/models/repo"
)

type ItemsService interface {
	Create(i *entity.Items) (*entity.Items, error)
	GetByID(id string) (*entity.Items, error)
	GetAll(limit, page int) ([]entity.Items, error)
	GetAllByType(limit, page int, itemType string) ([]entity.Items, error)
	Update(id string, i *entity.Items) (*entity.Items, error)
	Delete(id string) error
}

type itemsService struct{ repo repo.ItemsRepo }

func NewItemsService(r repo.ItemsRepo) ItemsService { return &itemsService{repo: r} }

func (s *itemsService) Create(i *entity.Items) (*entity.Items, error) {
	if i == nil {
		return nil, errors.New("item nil")
	}
	// trust repo to persist ItemType if provided
	if err := s.repo.Create(i); err != nil {
		return nil, err
	}
	return i, nil
}

func (s *itemsService) GetByID(id string) (*entity.Items, error) { return s.repo.GetByID(id) }

func (s *itemsService) GetAll(limit, page int) ([]entity.Items, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	list, err := s.repo.ListPage(limit, offset)
	if err != nil {
		return nil, err
	}
	out := make([]entity.Items, 0, len(list))
	for _, it := range list {
		out = append(out, *it)
	}
	return out, nil
}

func (s *itemsService) GetAllByType(limit, page int, itemType string) ([]entity.Items, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	list, err := s.repo.ListPageByType(limit, offset, itemType)
	if err != nil {
		return nil, err
	}
	out := make([]entity.Items, 0, len(list))
	for _, it := range list {
		out = append(out, *it)
	}
	return out, nil
}

func (s *itemsService) Update(id string, i *entity.Items) (*entity.Items, error) {
	if id == "" || i == nil {
		return nil, errors.New("invalid input")
	}
	i.IdItem = id
	if err := s.repo.Update(i); err != nil {
		return nil, err
	}
	return i, nil
}

func (s *itemsService) Delete(id string) error {
	if id == "" {
		return errors.New("id required")
	}
	return s.repo.Delete(id)
}
