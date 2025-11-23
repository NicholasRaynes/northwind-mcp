package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-api/internal/db"
	"github.com/nicholasraynes/northwind-api/internal/models"
)

// GET /suppliers
// Optional filters: ?country=Japan&company=Exotic
func GetSuppliers(c *gin.Context) {
	country := c.Query("country")
	company := c.Query("company")

	query := `
		SELECT supplier_id, company_name, contact_name, contact_title,
		       city, country, phone, fax, homepage
		FROM suppliers
	`
	conditions := []string{}
	args := []any{}

	if country != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(country) = LOWER($%d)", len(args)+1))
		args = append(args, country)
	}
	if company != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(company_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+company+"%")
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

	suppliers := []models.Supplier{}
	for rows.Next() {
		var s models.Supplier
		err := rows.Scan(
			&s.SupplierID,
			&s.CompanyName,
			&s.ContactName,
			&s.ContactTitle,
			&s.City,
			&s.Country,
			&s.Phone,
			&s.Fax,
			&s.HomePage,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		suppliers = append(suppliers, s)
	}

	c.JSON(http.StatusOK, gin.H{"count": len(suppliers), "data": suppliers})
}
