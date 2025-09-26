package repo

import (
	"faizalmaulana/lsp/models/entity"

	"gorm.io/gorm"
)

type PivotItemsToTransactionsRepo interface {
	BulkCreate(items []entity.PivotItemsToTransaction) error
	ListByTransaction(idTransaction string) ([]entity.PivotItemsToTransaction, error)
	DeleteByTransaction(idTransaction string) error
}

type GormPivotItemsToTransactionsRepo struct{ db *gorm.DB }

func NewGormPivotItemsToTransactionsRepo(db *gorm.DB) PivotItemsToTransactionsRepo {
	return &GormPivotItemsToTransactionsRepo{db: db}
}

func (r *GormPivotItemsToTransactionsRepo) BulkCreate(items []entity.PivotItemsToTransaction) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.Create(&items).Error
}

func (r *GormPivotItemsToTransactionsRepo) ListByTransaction(idTransaction string) ([]entity.PivotItemsToTransaction, error) {
	var out []entity.PivotItemsToTransaction
	if err := r.db.Where("id_transaction = ? AND is_deleted = ?", idTransaction, false).Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *GormPivotItemsToTransactionsRepo) DeleteByTransaction(idTransaction string) error {
	return r.db.Model(&entity.PivotItemsToTransaction{}).
		Where("id_transaction = ?", idTransaction).
		Update("is_deleted", true).Error
}
