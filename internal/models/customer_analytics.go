package models

type TopCustomer struct {
	CustomerID   string  `json:"customer_id" db:"customer_id"`
	CompanyName  string  `json:"company_name" db:"company_name"`
	Country      string  `json:"country" db:"country"`
	TotalSales   float64 `json:"total_sales" db:"total_sales"`
	OrderCount   int     `json:"order_count" db:"order_count"`
	AverageOrder float64 `json:"average_order" db:"average_order"`
}
