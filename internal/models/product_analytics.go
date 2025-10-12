package models

type TopProduct struct {
	ProductID    int     `json:"product_id" db:"product_id"`
	ProductName  string  `json:"product_name" db:"product_name"`
	CategoryName string  `json:"category_name" db:"category_name"`
	SupplierName string  `json:"supplier_name" db:"supplier_name"`
	UnitsSold    int     `json:"units_sold" db:"units_sold"`
	TotalRevenue float64 `json:"total_revenue" db:"total_revenue"`
	AveragePrice float64 `json:"average_price" db:"average_price"`
}
