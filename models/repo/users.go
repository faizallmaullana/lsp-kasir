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
	List(page, perPage int) ([]*entity.Users, int64, error)
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
	if err := r.db.First(&u, "id_user = ? AND is_deleted = ?", id, false).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *GormUsersRepo) GetByEmail(email string) (*entity.Users, error) {
	var u entity.Users
	if err := r.db.First(&u, "email = ? AND is_deleted = ?", email, false).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *GormUsersRepo) List(page, perPage int) ([]*entity.Users, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}

	var out []*entity.Users
	var total int64

	q := r.db.Model(&entity.Users{}).Where("is_deleted = ?", false)

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := q.Limit(perPage).Offset(offset).Find(&out).Error; err != nil {
		return nil, 0, err
	}

	totalPages := (total + int64(perPage) - 1) / int64(perPage)
	return out, totalPages, nil
}

func (r *GormUsersRepo) Update(u *entity.Users) error {
	if err := r.db.Model(&entity.Users{}).Where("id_user = ?", u.IdUser).Updates(u).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormUsersRepo) Delete(id string) error {
	if err := r.db.Model(&entity.Users{}).Where("id_user = ?", id).Update("is_deleted", true).Error; err != nil {
		return err
	}
	return nil
}
