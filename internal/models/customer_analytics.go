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

type CustomerLTV struct {
	CustomerID  string  `json:"customer_id" db:"customer_id"`
	CompanyName string  `json:"company_name" db:"company_name"`
	Country     string  `json:"country" db:"country"`
	TotalSales  float64 `json:"total_sales" db:"total_sales"`
	FirstOrder  string  `json:"first_order" db:"first_order"`
	LastOrder   string  `json:"last_order" db:"last_order"`
	OrderCount  int     `json:"order_count" db:"order_count"`
	AvgOrder    float64 `json:"avg_order" db:"avg_order"`
}

type CustomerRetention struct {
	CustomerID     string `json:"customer_id" db:"customer_id"`
	CompanyName    string `json:"company_name" db:"company_name"`
	Country        string `json:"country" db:"country"`
	FirstOrderYear int    `json:"first_order_year" db:"first_order_year"`
	LastOrderYear  int    `json:"last_order_year" db:"last_order_year"`
	OrderCount     int    `json:"order_count" db:"order_count"`
	ActiveYears    int    `json:"active_years" db:"active_years"`
	RepeatCustomer bool   `json:"repeat_customer" db:"repeat_customer"`
}
