package services

import (
	"context"
	"errors"
	"fmt"

	"faizalmaulana/lsp/models/entity"
	"faizalmaulana/lsp/models/repo"
)

type UsersService interface {
	Create(u *entity.Users) (*entity.Users, error)
	GetAll(count, page int) ([]entity.Users, error)
	GetByID(id string) (*entity.Users, error)
	GetByEmail(email string) (*entity.Users, error)
	Update(id string, u *entity.Users) (*entity.Users, error)
	Delete(id string) error
}

type usersService struct {
	users repo.UsersRepo
}

func NewUsersService(u repo.UsersRepo) UsersService {
	return &usersService{users: u}
}

func (s *usersService) Create(u *entity.Users) (*entity.Users, error) {
	if err := s.users.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *usersService) GetAll(count, page int) ([]entity.Users, error) {
	if count <= 0 {
		count = 10
	}
	if count > 100 {
		count = 100 
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * count

	list, err := s.users.ListPage(count, offset)
	if err != nil {
		return nil, err
	}
	out := make([]entity.Users, 0, len(list))
	for _, u := range list {
		out = append(out, *u)
	}
	return out, nil
}

func (s *usersService) GetByID(id string) (*entity.Users, error) {
	u, err := s.users.GetByID(id)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *usersService) GetByEmail(email string) (*entity.Users, error) {
	u, err := s.users.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *usersService) Update(id string, u *entity.Users) (*entity.Users, error) {
	if id == "" || u == nil {
		return nil, errors.New("invalid input")
	}
	u.IdUser = id
	if err := s.users.Update(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *usersService) Delete(id string) error {
	if id == "" {
		return errors.New("id required")
	}
	return s.users.Delete(id)
}

// optional
func (s *usersService) withContext(ctx context.Context) repo.UsersRepo {
	return s.users
}

func debugUser(u *entity.Users) {
	if u == nil {
		fmt.Println("debugUser: nil user")
		return
	}
	fmt.Printf("debugUser: id=%s email=%s profiles=%d\n", u.IdUser, u.Email, len(u.Profiles))
}
