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

type ShippingCosts struct {
	ShipperID      int     `json:"shipper_id" db:"shipper_id"`
	CompanyName    string  `json:"company_name" db:"company_name"`
	TotalOrders    int     `json:"total_orders" db:"total_orders"`
	TotalFreight   float64 `json:"total_freight" db:"total_freight"`
	AvgFreight     float64 `json:"avg_freight" db:"avg_freight"`
	TopDestination string  `json:"top_destination" db:"top_destination"`
	Year           *int    `json:"year,omitempty" db:"year"`
}

type DeliveryTimes struct {
	ShipperName     string  `json:"shipper_name" db:"shipper_name"`
	EmployeeName    string  `json:"employee_name" db:"employee_name"`
	TotalOrders     int     `json:"total_orders" db:"total_orders"`
	AvgDeliveryDays float64 `json:"avg_delivery_days" db:"avg_delivery_days"`
	MaxDeliveryDays float64 `json:"max_delivery_days" db:"max_delivery_days"`
	MinDeliveryDays float64 `json:"min_delivery_days" db:"min_delivery_days"`
	LateShipments   int     `json:"late_shipments" db:"late_shipments"`
	OnTimeRate      float64 `json:"on_time_rate" db:"on_time_rate"`
	Year            *int    `json:"year,omitempty" db:"year"`
}
