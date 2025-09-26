package dto

type CreateItemRequest struct {
	ItemName    string  `json:"item_name" binding:"required"`
	ItemType    string  `json:"item_type"`
	IsAvailable *bool   `json:"is_available"`
	Price       float64 `json:"price" binding:"required"`
	Description string  `json:"description"`
	ImageUrl    string  `json:"image_url"`
}

type UpdateItemRequest struct {
	ItemName    *string  `json:"item_name"`
	ItemType    *string  `json:"item_type"`
	IsAvailable *bool    `json:"is_available"`
	Price       *float64 `json:"price"`
	Description *string  `json:"description"`
	ImageUrl    *string  `json:"image_url"`
}
