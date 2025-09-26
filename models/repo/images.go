package repo

import (
	"errors"
	"faizalmaulana/lsp/models/entity"

	"gorm.io/gorm"
)

type ImagesRepo interface {
	Create(img *entity.Images) error
	GetByID(id string) (*entity.Images, error)
	Delete(id string) error
}

type gormImagesRepo struct{ db *gorm.DB }

func NewGormImagesRepo(db *gorm.DB) ImagesRepo { return &gormImagesRepo{db: db} }

func (r *gormImagesRepo) Create(img *entity.Images) error {
	if img == nil {
		return errors.New("nil image")
	}
	return r.db.Create(img).Error
}

func (r *gormImagesRepo) GetByID(id string) (*entity.Images, error) {
	var out entity.Images
	if err := r.db.Where("id_image = ? AND is_deleted = false", id).First(&out).Error; err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *gormImagesRepo) Delete(id string) error {
	return r.db.Model(&entity.Images{}).Where("id_image = ?", id).Update("is_deleted", true).Error
}
