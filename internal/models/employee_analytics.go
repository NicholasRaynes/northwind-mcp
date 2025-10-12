package models

type EmployeePerformance struct {
	EmployeeID   int     `json:"employee_id" db:"employee_id"`
	FullName     string  `json:"full_name" db:"full_name"`
	Title        string  `json:"title" db:"title"`
	Country      string  `json:"country" db:"country"`
	OrderCount   int     `json:"order_count" db:"order_count"`
	TotalRevenue float64 `json:"total_revenue" db:"total_revenue"`
	AvgOrder     float64 `json:"avg_order" db:"avg_order"`
}
