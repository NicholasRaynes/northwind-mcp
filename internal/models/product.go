package models

type Product struct {
	ProductID       int     `json:"product_id" db:"product_id"`
	ProductName     string  `json:"product_name" db:"product_name"`
	SupplierID      int     `json:"supplier_id" db:"supplier_id"`
	SupplierName    string  `json:"supplier_name" db:"supplier_name"`
	CategoryID      int     `json:"category_id" db:"category_id"`
	CategoryName    string  `json:"category_name" db:"category_name"`
	QuantityPerUnit string  `json:"quantity_per_unit" db:"quantity_per_unit"`
	UnitPrice       float64 `json:"unit_price" db:"unit_price"`
	UnitsInStock    int     `json:"units_in_stock" db:"units_in_stock"`
	Discontinued    bool    `json:"discontinued" db:"discontinued"`
}
