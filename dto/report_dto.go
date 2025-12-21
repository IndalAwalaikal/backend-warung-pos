package dto

type DailyReportResponse struct {
	Date              string      `json:"date"`
	TotalRevenue      float64     `json:"total_revenue"`
	TotalTransactions int         `json:"total_transactions"`
	TotalItems        int         `json:"total_items"`
	BestSellers       interface{} `json:"best_sellers"`
	Transactions      interface{} `json:"transactions"`
}
