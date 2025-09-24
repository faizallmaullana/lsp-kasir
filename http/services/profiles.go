package services

import (
	"errors"

	"faizalmaulana/lsp/models/entity"
	"faizalmaulana/lsp/models/repo"
)

type ProfilesService interface {
	Create(p *entity.Profiles) (*entity.Profiles, error)
	GetAll() ([]entity.Profiles, error)
	GetByID(id string) (*entity.Profiles, error)
	Update(id string, p *entity.Profiles) (*entity.Profiles, error)
	Delete(id string) error
}

type profilesService struct {
	profile repo.ProfilesRepo
}

func NewProfilesService(u repo.ProfilesRepo) ProfilesService {
	return &profilesService{profile: u}
}

func (s *profilesService) Create(p *entity.Profiles) (*entity.Profiles, error) {
	if p == nil {
		return nil, errors.New("profile is nil")
	}
	if err := s.profile.Create(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *profilesService) GetAll() ([]entity.Profiles, error) {
	list, err := s.profile.List()
	if err != nil {
		return nil, err
	}
	out := make([]entity.Profiles, 0, len(list))
	for _, p := range list {
		out = append(out, *p)
	}
	return out, nil
}

func (s *profilesService) GetByID(id string) (*entity.Profiles, error) {
	if id == "" {
		return nil, errors.New("id required")
	}
	return s.profile.GetByID(id)
}

func (s *profilesService) Update(id string, p *entity.Profiles) (*entity.Profiles, error) {
	if id == "" || p == nil {
		return nil, errors.New("invalid input")
	}
	p.IdProfile = id
	if err := s.profile.Update(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *profilesService) Delete(id string) error {
	if id == "" {
		return errors.New("id required")
	}
	return s.profile.Delete(id)
}
