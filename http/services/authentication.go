package services

import (
	"errors"

	"faizalmaulana/lsp/models/entity"
	"faizalmaulana/lsp/models/repo"

	"golang.org/x/crypto/bcrypt"
)

type AuthenticationService interface {
	Login(email, password string) (*entity.Users, error)
}

type authenticationService struct {
	users repo.UsersRepo
}

func NewAuthenticationService(u repo.UsersRepo) AuthenticationService {
	return &authenticationService{users: u}
}

func (s *authenticationService) Login(email, password string) (*entity.Users, error) {
	user, err := s.users.GetByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
