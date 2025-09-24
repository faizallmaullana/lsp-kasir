package repo

import (
	"errors"
	"faizalmaulana/lsp/models/entity"

	"gorm.io/gorm"
)

type TransactionsRepo interface {
	Create(u *entity.Transactions) error
	GetByID(id string) (*entity.Transactions, error)
	List() ([]*entity.Transactions, error)
	Update(u *entity.Transactions) error
	Delete(id string) error
}

type GormTransactionsRepo struct {
	db *gorm.DB
}

func NewGormTransactionsRepo(db *gorm.DB) TransactionsRepo {
	return &GormTransactionsRepo{db: db}
}

func (r *GormTransactionsRepo) Create(u *entity.Transactions) error {
	return r.db.Create(u).Error
}

func (r *GormTransactionsRepo) GetByID(id string) (*entity.Transactions, error) {
	var u entity.Transactions
	if err := r.db.First(&u, "id_transaction = ? AND is_deleted = ?", id, false).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *GormTransactionsRepo) List() ([]*entity.Transactions, error) {
	var out []*entity.Transactions
	if err := r.db.Where("is_deleted = ?", false).Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *GormTransactionsRepo) Update(u *entity.Transactions) error {
	if err := r.db.Model(&entity.Transactions{}).Where("id_transaction = ?", u.IdTransaction).Updates(u).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormTransactionsRepo) Delete(id string) error {
	if err := r.db.Model(&entity.Transactions{}).Where("id_transaction = ?", id).Update("is_deleted", true).Error; err != nil {
		return err
	}
	return nil
}
