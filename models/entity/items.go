package entity

import "time"

type Items struct {
	IdItem string `json:"id_item" gorm:"type:varchar(36);unique;primaryKey;not null"`

	ItemName    string  `json:"item_name" gorm:"type:varchar(255);not null"`
	IsAvailable bool    `json:"is_available" gorm:"type:boolean;default:true"`
	Price       float64 `json:"price" gorm:"type:decimal(10,2);not null"`
	Description string  `json:"description" gorm:"type:text"`
	ImageUrl    string  `json:"image_url" gorm:"type:varchar(255)"`

	Timestamp time.Time `json:"timestamp" gorm:"autoCreateTime"`
	IsDeleted bool      `json:"is_deleted" gorm:"type:boolean;default:false"`

	PivotItemsToTransaction []PivotItemsToTransaction `json:"pivot_items_to_transaction" gorm:"foreignKey:IdItem;references:IdItem;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
