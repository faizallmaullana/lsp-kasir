package services

import (
	"faizalmaulana/lsp/helper"
	"faizalmaulana/lsp/models/entity"
	"faizalmaulana/lsp/models/repo"
)

type SessionService interface {
	GetByUserID(userID string) (*entity.Sessions, error)
	Create(idUser string) (*entity.Sessions, error)
	GetAll(limit, page int) ([]entity.Sessions, error)
}

type sessionService struct {
	users repo.SessionsRepo
}

func NewSessionService(u repo.SessionsRepo) SessionService {
	return &sessionService{users: u}
}

func (s *sessionService) GetByUserID(userID string) (*entity.Sessions, error) {
	return s.users.GetByIdUser(userID)
}

func (s *sessionService) Create(idUser string) (*entity.Sessions, error) {
	session := &entity.Sessions{
		IdSession: helper.Uuid(),
		IdUser:    idUser,
	}

	if err := s.users.Create(session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *sessionService) GetAll(limit, page int) ([]entity.Sessions, error) {
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
	list, err := s.users.ListPage(limit, offset)
	if err != nil {
		return nil, err
	}
	out := make([]entity.Sessions, 0, len(list))
	for _, v := range list {
		out = append(out, *v)
	}
	return out, nil
}
