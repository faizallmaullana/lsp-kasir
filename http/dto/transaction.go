package dto

type TransactionItemRequest struct {
	IdItem   string `json:"id_item" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,min=1"`
}

type CreateTransactionRequest struct {
	BuyerContact string                   `json:"buyer_contact"`
	Items        []TransactionItemRequest `json:"items" binding:"required"`
}

type UpdateTransactionRequest struct {
	BuyerContact *string `json:"buyer_contact"`
}

type TransactionItemDetail struct {
	IdItem   string  `json:"id_item"`
	ItemName string  `json:"item_name"`
	ImageUrl string  `json:"image_url"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}
