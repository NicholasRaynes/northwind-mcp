package models

type TopCustomer struct {
	CustomerID   string  `json:"customer_id" db:"customer_id"`
	CompanyName  string  `json:"company_name" db:"company_name"`
	Country      string  `json:"country" db:"country"`
	TotalSales   float64 `json:"total_sales" db:"total_sales"`
	OrderCount   int     `json:"order_count" db:"order_count"`
	AverageOrder float64 `json:"average_order" db:"average_order"`
}

type CustomerOrderSummary struct {
	CustomerID  string  `json:"customer_id" db:"customer_id"`
	CompanyName string  `json:"company_name" db:"company_name"`
	OrderID     int     `json:"order_id" db:"order_id"`
	OrderDate   string  `json:"order_date" db:"order_date"`
	TotalAmount float64 `json:"total_amount" db:"total_amount"`
	ShippedDate string  `json:"shipped_date" db:"shipped_date"`
	Country     string  `json:"country" db:"country"`
}
