package entity

import "time"

type Transactions struct {
	IdTransaction string `json:"id_transaction" gorm:"type:varchar(36);unique;primaryKey;not null"`
	IdUser        string `json:"id_user" gorm:"type:varchar(36);not null;index"`

	BuyerContact string  `json:"buyer_contact" gorm:"type:varchar(120)"`
	TotalPrice   float64 `json:"total_price"`

	IsDeleted bool `json:"is_deleted" gorm:"type:boolean;default:false"`

	Timestamp time.Time `json:"timestamp" gorm:"autoCreateTime"`

	PivotItemsToTransaction []PivotItemsToTransaction `json:"pivot_items_to_transaction" gorm:"foreignKey:IdTransaction;references:IdTransaction;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
