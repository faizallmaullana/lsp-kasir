package dto

type CreateItemRequest struct {
	ItemName    string  `json:"item_name" binding:"required"`
	ItemType    string  `json:"item_type"`
	IsAvailable *bool   `json:"is_available"`
	Price       float64 `json:"price" binding:"required"`
	Description string  `json:"description"`
	ImageUrl    string  `json:"image_url"`
	ImageBase64 string  `json:"image_base64"` // optional: when provided, will be saved to storages/images and file name stored
	ImageType   string  `json:"image_type"`   // optional: MIME type like image/png when using base64
}

type UpdateItemRequest struct {
	ItemName    *string  `json:"item_name"`
	ItemType    *string  `json:"item_type"`
	IsAvailable *bool    `json:"is_available"`
	Price       *float64 `json:"price"`
	Description *string  `json:"description"`
	ImageUrl    *string  `json:"image_url"`
	ImageBase64 *string  `json:"image_base64"`
	ImageType   *string  `json:"image_type"`
}
