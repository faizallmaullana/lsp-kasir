package repo

import (
	"errors"
	"faizalmaulana/lsp/models/entity"

	"gorm.io/gorm"
)

type SessionsRepo interface {
	Create(u *entity.Sessions) error
	GetByID(id string) (*entity.Sessions, error)
	GetByIdUser(id string) (*entity.Sessions, error)
	List() ([]*entity.Sessions, error)
	ListPage(limit, offset int) ([]*entity.Sessions, error)
	Update(u *entity.Sessions) error
	Delete(id string) error
}

type GormSessionsRepo struct {
	db *gorm.DB
}

func NewGormSessionsRepo(db *gorm.DB) SessionsRepo {
	return &GormSessionsRepo{db: db}
}

func (r *GormSessionsRepo) Create(u *entity.Sessions) error {
	return r.db.Create(u).Error
}

func (r *GormSessionsRepo) GetByID(id string) (*entity.Sessions, error) {
	var u entity.Sessions
	if err := r.db.First(&u, "id_session = ? AND is_deleted = ?", id, false).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *GormSessionsRepo) GetByIdUser(id string) (*entity.Sessions, error) {
	var u entity.Sessions
	if err := r.db.First(&u, "id_user = ? AND is_deleted = ?", id, false).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *GormSessionsRepo) List() ([]*entity.Sessions, error) {
	var out []*entity.Sessions
	if err := r.db.Where("is_deleted = ?", false).Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *GormSessionsRepo) ListPage(limit, offset int) ([]*entity.Sessions, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	var out []*entity.Sessions
	if err := r.db.Where("is_deleted = ?", false).Limit(limit).Offset(offset).Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *GormSessionsRepo) Update(u *entity.Sessions) error {
	if err := r.db.Model(&entity.Sessions{}).Where("id_session = ?", u.IdSession).Updates(u).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormSessionsRepo) Delete(id string) error {
	if err := r.db.Model(&entity.Sessions{}).Where("id_session = ?", id).Update("is_deleted", true).Error; err != nil {
		return err
	}
	return nil
}
