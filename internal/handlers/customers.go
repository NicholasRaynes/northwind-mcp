package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-mcp/internal/db"
	"github.com/nicholasraynes/northwind-mcp/internal/models"
)

// GET /customers
func GetCustomers(c *gin.Context) {
	rows, err := db.DB.Query(`
		SELECT customer_id, company_name, contact_name, contact_title, address,
		       city, region, postal_code, country, phone, fax
		FROM customers
		ORDER BY company_name
	`)
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

	c.JSON(http.StatusOK, gin.H{"data": customers})
}
