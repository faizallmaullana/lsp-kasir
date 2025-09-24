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
	GetAll() ([]entity.Users, error)
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

// Create persists a new user.
func (s *usersService) Create(u *entity.Users) (*entity.Users, error) {
	if err := s.users.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

// GetAll returns a list of users.
func (s *usersService) GetAll() ([]entity.Users, error) {
	list, err := s.users.List()
	if err != nil {
		return nil, err
	}

	// Simple pagination: return first page by default.
	// TODO: accept page/limit parameters (via function args, context, or caller config).
	const (
		defaultPage  = 1
		defaultLimit = 10
	)

	page := defaultPage
	limit := defaultLimit

	// calculate slice bounds
	start := (page - 1) * limit
	if start >= len(list) {
		// no results on this page
		return []entity.Users{}, nil
	}
	end := start + limit
	if end > len(list) {
		end = len(list)
	}

	// convert []*entity.Users to []entity.Users for interface contract
	out := make([]entity.Users, 0, end-start)
	for _, p := range list[start:end] {
		out = append(out, *p)
	}
	return out, nil
}

// GetByID returns a user by id (with basic related data if already preloaded by repo).
func (s *usersService) GetByID(id string) (*entity.Users, error) {
	u, err := s.users.GetByID(id)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// GetByEmail returns a user by email.
func (s *usersService) GetByEmail(email string) (*entity.Users, error) {
	u, err := s.users.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Update modifies an existing user.
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

// Delete removes a user.
func (s *usersService) Delete(id string) error {
	if id == "" {
		return errors.New("id required")
	}
	return s.users.Delete(id)
}

// Optional: example of method with context (future extension)
func (s *usersService) withContext(ctx context.Context) repo.UsersRepo {
	// placeholder for future repo methods that accept context
	return s.users
}

// Debug helper (not exported)
func debugUser(u *entity.Users) {
	if u == nil {
		fmt.Println("debugUser: nil user")
		return
	}
	fmt.Printf("debugUser: id=%s email=%s profiles=%d\n", u.IdUser, u.Email, len(u.Profiles))
}
