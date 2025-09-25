package dto

// CreateItemRequest represents payload to create an item.
type CreateItemRequest struct {
	ItemName    string  `json:"item_name" binding:"required"`
	IsAvailable *bool   `json:"is_available"`
	Price       float64 `json:"price" binding:"required"`
	Description string  `json:"description"`
	ImageUrl    string  `json:"image_url"`
}

// UpdateItemRequest represents payload to update an item.
type UpdateItemRequest struct {
	ItemName    *string  `json:"item_name"`
	IsAvailable *bool    `json:"is_available"`
	Price       *float64 `json:"price"`
	Description *string  `json:"description"`
	ImageUrl    *string  `json:"image_url"`
}
