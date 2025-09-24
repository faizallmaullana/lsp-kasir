package services

import (
    "errors"
    "faizalmaulana/lsp/models/entity"
    "faizalmaulana/lsp/models/repo"
)

type TransactionsService interface {
    Create(t *entity.Transactions) (*entity.Transactions, error)
    GetByID(id string) (*entity.Transactions, error)
    GetAll(limit, page int) ([]entity.Transactions, error)
    Update(id string, t *entity.Transactions) (*entity.Transactions, error)
    Delete(id string) error
}

type transactionsService struct { repo repo.TransactionsRepo }

func NewTransactionsService(r repo.TransactionsRepo) TransactionsService { return &transactionsService{repo: r} }

func (s *transactionsService) Create(t *entity.Transactions) (*entity.Transactions, error) {
    if t == nil { return nil, errors.New("transaction nil") }
    if err := s.repo.Create(t); err != nil { return nil, err }
    return t, nil
}

func (s *transactionsService) GetByID(id string) (*entity.Transactions, error) { return s.repo.GetByID(id) }

func (s *transactionsService) GetAll(limit, page int) ([]entity.Transactions, error) {
    if limit <= 0 { limit = 10 }
    if limit > 100 { limit = 100 }
    if page <= 0 { page = 1 }
    offset := (page - 1) * limit
    list, err := s.repo.ListPage(limit, offset)
    if err != nil { return nil, err }
    out := make([]entity.Transactions, 0, len(list))
    for _, t := range list { out = append(out, *t) }
    return out, nil
}

func (s *transactionsService) Update(id string, t *entity.Transactions) (*entity.Transactions, error) {
    if id == "" || t == nil { return nil, errors.New("invalid input") }
    t.IdTransaction = id
    if err := s.repo.Update(t); err != nil { return nil, err }
    return t, nil
}

func (s *transactionsService) Delete(id string) error { if id == "" { return errors.New("id required") }; return s.repo.Delete(id) }
