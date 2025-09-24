package repo

import (
	"errors"
	"faizalmaulana/lsp/models/entity"

	"gorm.io/gorm"
)

type ItemsRepo interface {
	Create(u *entity.Items) error
	GetByID(id string) (*entity.Items, error)
	List() ([]*entity.Items, error)
	Update(u *entity.Items) error
	Delete(id string) error
}

type GormItemsRepo struct {
	db *gorm.DB
}

func NewGormItemsRepo(db *gorm.DB) ItemsRepo {
	return &GormItemsRepo{db: db}
}

func (r *GormItemsRepo) Create(u *entity.Items) error {
	return r.db.Create(u).Error
}

func (r *GormItemsRepo) GetByID(id string) (*entity.Items, error) {
	var u entity.Items
	if err := r.db.Where("id_item = ? AND is_deleted = ?", id, false).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *GormItemsRepo) List() ([]*entity.Items, error) {
	var out []*entity.Items
	if err := r.db.Where("is_deleted = ?", false).Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *GormItemsRepo) Update(u *entity.Items) error {
	if err := r.db.Model(&entity.Items{}).Where("id_item = ?", u.IdItem).Updates(u).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormItemsRepo) Delete(id string) error {
	if err := r.db.Model(&entity.Items{}).Where("id_item = ?", id).Update("isdeleted", true).Error; err != nil {
		return err
	}
	return nil
}
