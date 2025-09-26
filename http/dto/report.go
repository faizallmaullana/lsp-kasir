package dto

type ReportResponse struct {
	Month int                 `json:"month"`
	Year  int                 `json:"year"`
	Total int                 `json:"total_transactions"`
	Sum   float64             `json:"sum_total_price"`
	Items []ReportTransaction `json:"transactions"`
}

type ReportTransaction struct {
	IdTransaction string  `json:"id_transaction"`
	TotalPrice    float64 `json:"total_price"`
	BuyerContact  string  `json:"buyer_contact"`
	Timestamp     string  `json:"timestamp"`
}

type TodayReportResponse struct {
	Date  string              `json:"date"`
	Total int                 `json:"total_transactions"`
	Sum   float64             `json:"sum_total_price"`
	Items []ReportTransaction `json:"transactions"`
}

type TodaySummaryResponse struct {
	Date              string  `json:"date"`
	TotalTransactions int     `json:"total_transactions"`
	TotalProductsSold int     `json:"total_products_sold"`
	SumTotalPrice     float64 `json:"sum_total_price"`
}
