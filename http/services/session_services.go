package services

import (
	"faizalmaulana/lsp/helper"
	"faizalmaulana/lsp/models/entity"
	"faizalmaulana/lsp/models/repo"
)

type SessionService interface {
	GetByUserID(userID string) (*entity.Sessions, error)
	Create(idUser string) (*entity.Sessions, error)
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
