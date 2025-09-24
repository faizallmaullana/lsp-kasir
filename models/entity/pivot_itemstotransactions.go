package entity

type PivotItemsToTransaction struct {
	IdTransaction string `json:"id_transaction" gorm:"type:varchar(36);not null;index"`
	IdItem        string `json:"id_item" gorm:"type:varchar(36);not null;index"`

	IsDeleted bool      `json:"is_deleted" gorm:"type:boolean;default:false"`

	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}
