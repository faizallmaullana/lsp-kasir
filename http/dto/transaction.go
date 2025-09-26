package dto

// TransactionItemRequest represents one purchased item in a transaction.
type TransactionItemRequest struct {
	IdItem   string `json:"id_item" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,min=1"`
}

// CreateTransactionRequest is the payload to create a cashier transaction.
type CreateTransactionRequest struct {
	BuyerContact string                   `json:"buyer_contact"`
	Items        []TransactionItemRequest `json:"items" binding:"required"`
}

// UpdateTransactionRequest allows updating mutable fields of a transaction.
type UpdateTransactionRequest struct {
	BuyerContact *string `json:"buyer_contact"`
}

// TransactionItemDetail is used in GET transaction by id response to include item info.
type TransactionItemDetail struct {
	IdItem   string  `json:"id_item"`
	ItemName string  `json:"item_name"`
	ImageUrl string  `json:"image_url"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}
