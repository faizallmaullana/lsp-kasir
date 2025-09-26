package dto

type ReportResponse struct {
	Month             int                 `json:"month"`
	Year              int                 `json:"year"`
	Total             int                 `json:"total_transactions"`
	Sum               float64             `json:"sum_total_price"`
	TotalProductsSold int                 `json:"total_products_sold"`
	AverageOrderValue float64             `json:"average_order_value"`
	MinOrderValue     float64             `json:"min_order_value"`
	MaxOrderValue     float64             `json:"max_order_value"`
	AvgItemsPerTx     float64             `json:"average_items_per_transaction"`
	TopItems          []TopItem           `json:"top_items"`
	Items             []ReportTransaction `json:"transactions"`
}

type ReportTransaction struct {
	IdTransaction string  `json:"id_transaction"`
	TotalPrice    float64 `json:"total_price"`
	BuyerContact  string  `json:"buyer_contact"`
	Timestamp     string  `json:"timestamp"`
}

type TodayReportResponse struct {
	Date              string              `json:"date"`
	Total             int                 `json:"total_transactions"`
	Sum               float64             `json:"sum_total_price"`
	TotalProductsSold int                 `json:"total_products_sold"`
	AverageOrderValue float64             `json:"average_order_value"`
	MinOrderValue     float64             `json:"min_order_value"`
	MaxOrderValue     float64             `json:"max_order_value"`
	AvgItemsPerTx     float64             `json:"average_items_per_transaction"`
	TopItems          []TopItem           `json:"top_items"`
	Items             []ReportTransaction `json:"transactions"`
}

type TodaySummaryResponse struct {
	Date              string    `json:"date"`
	TotalTransactions int       `json:"total_transactions"`
	TotalProductsSold int       `json:"total_products_sold"`
	SumTotalPrice     float64   `json:"sum_total_price"`
	AverageOrderValue float64   `json:"average_order_value"`
	MinOrderValue     float64   `json:"min_order_value"`
	MaxOrderValue     float64   `json:"max_order_value"`
	AvgItemsPerTx     float64   `json:"average_items_per_transaction"`
	TopItems          []TopItem `json:"top_items"`
}

type TopItem struct {
	IdItem       string  `json:"id_item"`
	ItemName     string  `json:"item_name"`
	ImageUrl     string  `json:"image_url"`
	QuantitySold int     `json:"quantity_sold"`
	Revenue      float64 `json:"revenue"`
}
