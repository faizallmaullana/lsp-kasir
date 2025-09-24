package repo

import (
	"errors"
	"faizalmaulana/lsp/models/entity"

	"gorm.io/gorm"
)

type ProfilesRepo interface {
	Create(u *entity.Profiles) error
	GetByID(id string) (*entity.Profiles, error)
	List() ([]*entity.Profiles, error)
	Update(u *entity.Profiles) error
	Delete(id string) error
}

type GormProfilesRepo struct {
	db *gorm.DB
}

func NewGormProfilesRepo(db *gorm.DB) ProfilesRepo {
	return &GormProfilesRepo{db: db}
}

func (r *GormProfilesRepo) Create(u *entity.Profiles) error {
	return r.db.Create(u).Error
}

func (r *GormProfilesRepo) GetByID(id string) (*entity.Profiles, error) {
	var u entity.Profiles
	if err := r.db.First(&u, "id_profile = ? AND is_deleted = ?", id, false).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *GormProfilesRepo) List() ([]*entity.Profiles, error) {
	var out []*entity.Profiles
	if err := r.db.Where("is_deleted = ?", false).Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *GormProfilesRepo) Update(u *entity.Profiles) error {
	if err := r.db.Model(&entity.Profiles{}).Where("id_profile = ?", u.IdProfile).Updates(u).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormProfilesRepo) Delete(id string) error {
	if err := r.db.Model(&entity.Profiles{}).Where("id_profile = ?", id).Update("is_deleted", true).Error; err != nil {
		return err
	}
	return nil
}
