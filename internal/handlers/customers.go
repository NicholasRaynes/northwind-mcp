package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-api/internal/db"
	"github.com/nicholasraynes/northwind-api/internal/models"
)

// GET /customers
// Optional parameters: country, city, id, customer_id, company_name, contact_name, contact_title, address, region, postal_code, phone, fax
func GetCustomers(c *gin.Context) {
	country := c.Query("country")
	city := c.Query("city")
	id := c.Query("id")
	customerID := c.Query("customer_id")
	companyName := c.Query("company_name")
	contactName := c.Query("contact_name")
	contactTitle := c.Query("contact_title")
	address := c.Query("address")
	region := c.Query("region")
	postalCode := c.Query("postal_code")
	phone := c.Query("phone")
	fax := c.Query("fax")

	// Base query
	query := `
		SELECT customer_id, company_name, contact_name, contact_title, address,
		       city, region, postal_code, country, phone, fax
		FROM customers
	`
	conditions := []string{}
	args := []any{}

	if id != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(customer_id) = LOWER($%d)", len(args)+1))
		args = append(args, id)
	}
	if country != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(country) = LOWER($%d)", len(args)+1))
		args = append(args, country)
	}
	if city != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(city) = LOWER($%d)", len(args)+1))
		args = append(args, city)
	}
	if customerID != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(customer_id) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+customerID+"%")
	}
	if companyName != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(company_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+companyName+"%")
	}
	if contactName != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(contact_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+contactName+"%")
	}
	if contactTitle != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(contact_title) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+contactTitle+"%")
	}
	if address != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(address) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+address+"%")
	}
	if region != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(region) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+region+"%")
	}
	if postalCode != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(postal_code) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+postalCode+"%")
	}
	if phone != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(phone) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+phone+"%")
	}
	if fax != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(fax) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+fax+"%")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY company_name"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	customers := []models.Customer{}

	for rows.Next() {
		var cust models.Customer
		err := rows.Scan(
			&cust.CustomerID,
			&cust.CompanyName,
			&cust.ContactName,
			&cust.ContactTitle,
			&cust.Address,
			&cust.City,
			&cust.Region,
			&cust.PostalCode,
			&cust.Country,
			&cust.Phone,
			&cust.Fax,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		customers = append(customers, cust)
	}

	filters := gin.H{}
	if id != "" {
		filters["id"] = id
	}
	if country != "" {
		filters["country"] = country
	}
	if city != "" {
		filters["city"] = city
	}
	if customerID != "" {
		filters["customer_id"] = customerID
	}
	if companyName != "" {
		filters["company_name"] = companyName
	}
	if contactName != "" {
		filters["contact_name"] = contactName
	}
	if contactTitle != "" {
		filters["contact_title"] = contactTitle
	}
	if address != "" {
		filters["address"] = address
	}
	if region != "" {
		filters["region"] = region
	}
	if postalCode != "" {
		filters["postal_code"] = postalCode
	}
	if phone != "" {
		filters["phone"] = phone
	}
	if fax != "" {
		filters["fax"] = fax
	}

	c.JSON(http.StatusOK, gin.H{
		"filters": filters,
		"count":   len(customers),
		"data":    customers,
	})
}
