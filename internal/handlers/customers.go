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
// Optional filters: ?country=Brazil&city=London&id=ALFKI
func GetCustomers(c *gin.Context) {
	country := c.Query("country")
	city := c.Query("city")
	id := c.Query("id")

	// Base query
	query := `
		SELECT customer_id, company_name, contact_name, contact_title, address,
		       city, region, postal_code, country, phone, fax
		FROM customers
	`
	conditions := []string{}
	args := []any{}

	// Apply filters dynamically
	if id != "" {
		conditions = append(conditions, "LOWER(customer_id) = LOWER($1)")
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

	c.JSON(http.StatusOK, gin.H{
		"count": len(customers),
		"data":  customers,
	})
}
