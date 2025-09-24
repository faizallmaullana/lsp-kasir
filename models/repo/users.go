package repo

import (
	"errors"

	"faizalmaulana/lsp/models/entity"

	"gorm.io/gorm"
)

// UsersRepo defines CRUD for users
type UsersRepo interface {
	Create(u *entity.Users) error
	GetByID(id string) (*entity.Users, error)
	GetByEmail(email string) (*entity.Users, error)
	List() ([]*entity.Users, error)
	ListPage(limit, offset int) ([]*entity.Users, error) // new: paginated list
	Update(u *entity.Users) error
	Delete(id string) error
}

// GormUsersRepo is a GORM implementation of UsersRepo
type GormUsersRepo struct {
	db *gorm.DB
}

func NewGormUsersRepo(db *gorm.DB) UsersRepo {
	return &GormUsersRepo{db: db}
}

func (r *GormUsersRepo) Create(u *entity.Users) error {
	return r.db.Create(u).Error
}

func (r *GormUsersRepo) GetByID(id string) (*entity.Users, error) {
	var u entity.Users
	if err := r.db.First(&u, "id_user = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *GormUsersRepo) GetByEmail(email string) (*entity.Users, error) {
	var u entity.Users
	if err := r.db.First(&u, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *GormUsersRepo) List() ([]*entity.Users, error) {
	var out []*entity.Users
	if err := r.db.Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *GormUsersRepo) ListPage(limit, offset int) ([]*entity.Users, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	var out []*entity.Users
	if err := r.db.Limit(limit).Offset(offset).Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *GormUsersRepo) Update(u *entity.Users) error {
	return r.db.Model(&entity.Users{}).Where("id_user = ?", u.IdUser).Updates(u).Error
}

func (r *GormUsersRepo) Delete(id string) error {
	return r.db.Delete(&entity.Users{}, "id_user = ?", id).Error
}
