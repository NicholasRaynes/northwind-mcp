package models

type SalesSummary struct {
	GroupKey   string  `json:"group_key" db:"group_key"` // Value dependent on the summary type
	TotalSales float64 `json:"total_sales" db:"total_sales"`
	OrderCount int     `json:"order_count" db:"order_count"`
}
