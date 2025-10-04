package models

type OrderDetail struct {
	OrderID       int     `json:"order_id" db:"order_id"`
	ProductID     int     `json:"product_id" db:"product_id"`
	ProductName   string  `json:"product_name" db:"product_name"`
	Category      string  `json:"category_name" db:"category_name"`
	Supplier      string  `json:"supplier_name" db:"supplier_name"`
	UnitPrice     float64 `json:"unit_price" db:"unit_price"`
	Quantity      int     `json:"quantity" db:"quantity"`
	Discount      float64 `json:"discount" db:"discount"`
	ExtendedPrice float64 `json:"extended_price" db:"extended_price"`
}
