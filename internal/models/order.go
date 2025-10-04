package models

import "time"

type Order struct {
	OrderID      int        `json:"order_id" db:"order_id"`
	CustomerID   string     `json:"customer_id" db:"customer_id"`
	CustomerName string     `json:"customer_name" db:"customer_name"`
	EmployeeName string     `json:"employee_name" db:"employee_name"`
	OrderDate    *time.Time `json:"order_date" db:"order_date"`
	RequiredDate *time.Time `json:"required_date" db:"required_date"`
	ShippedDate  *time.Time `json:"shipped_date" db:"shipped_date"`
	ShipVia      int        `json:"ship_via" db:"ship_via"`
	ShipperName  string     `json:"shipper_name" db:"shipper_name"`
	Freight      float64    `json:"freight" db:"freight"`
	ShipName     string     `json:"ship_name" db:"ship_name"`
	ShipAddress  string     `json:"ship_address" db:"ship_address"`
	ShipCity     string     `json:"ship_city" db:"ship_city"`
	ShipRegion   *string    `json:"ship_region" db:"ship_region"`
	ShipPostal   *string    `json:"ship_postal" db:"ship_postal_code"`
	ShipCountry  string     `json:"ship_country" db:"ship_country"`
}
