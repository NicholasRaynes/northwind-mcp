package models

type Supplier struct {
	SupplierID   int     `json:"supplier_id" db:"supplier_id"`
	CompanyName  string  `json:"company_name" db:"company_name"`
	ContactName  string  `json:"contact_name" db:"contact_name"`
	ContactTitle string  `json:"contact_title" db:"contact_title"`
	City         string  `json:"city" db:"city"`
	Country      string  `json:"country" db:"country"`
	Phone        string  `json:"phone" db:"phone"`
	Fax          *string `json:"fax" db:"fax"`
	HomePage     *string `json:"homepage" db:"homepage"`
}
