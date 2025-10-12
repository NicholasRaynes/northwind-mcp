package models

type SalesSummary struct {
	GroupKey   string  `json:"group_key" db:"group_key"` // e.g., Country
	TotalSales float64 `json:"total_sales" db:"total_sales"`
	OrderCount int     `json:"order_count" db:"order_count"`
}
