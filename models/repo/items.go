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
	ListPage(limit, offset int) ([]*entity.Items, error)
	ListPageByType(limit, offset int, itemType string) ([]*entity.Items, error)
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

func (r *GormItemsRepo) ListPage(limit, offset int) ([]*entity.Items, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	var out []*entity.Items
	if err := r.db.Where("is_deleted = ?", false).Limit(limit).Offset(offset).Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *GormItemsRepo) ListPageByType(limit, offset int, itemType string) ([]*entity.Items, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	var out []*entity.Items

	query := r.db.Where("is_deleted = ?", false)
	if itemType != "" {
		query = query.Where("item_type = ?", itemType)
	}

	if err := query.Limit(limit).Offset(offset).Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *GormItemsRepo) Update(u *entity.Items) error {
	return r.db.Model(&entity.Items{}).Where("id_item = ?", u.IdItem).Updates(u).Error
}

func (r *GormItemsRepo) Delete(id string) error {
	return r.db.Model(&entity.Items{}).Where("id_item = ?", id).Update("is_deleted", true).Error
}
