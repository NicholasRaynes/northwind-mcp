package models

type Customer struct {
	CustomerID   string  `json:"customer_id" db:"customer_id"`
	CompanyName  string  `json:"company_name" db:"company_name"`
	ContactName  string  `json:"contact_name" db:"contact_name"`
	ContactTitle string  `json:"contact_title" db:"contact_title"`
	Address      string  `json:"address" db:"address"`
	City         string  `json:"city" db:"city"`
	Region       *string `json:"region" db:"region"`
	PostalCode   *string `json:"postal_code" db:"postal_code"`
	Country      string  `json:"country" db:"country"`
	Phone        string  `json:"phone" db:"phone"`
	Fax          *string `json:"fax" db:"fax"`
}
