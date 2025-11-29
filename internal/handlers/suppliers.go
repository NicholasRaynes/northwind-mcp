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
// Optional parameters: country, supplier_id, company_name, contact_name, contact_title, city, phone, fax
func GetSuppliers(c *gin.Context) {
	country := c.Query("country")
	supplierID := c.Query("supplier_id")
	companyName := c.Query("company_name")
	contactName := c.Query("contact_name")
	contactTitle := c.Query("contact_title")
	city := c.Query("city")
	phone := c.Query("phone")
	fax := c.Query("fax")

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
	if supplierID != "" {
		conditions = append(conditions, fmt.Sprintf("CAST(supplier_id AS TEXT) LIKE $%d", len(args)+1))
		args = append(args, "%"+supplierID+"%")
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
	if city != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(city) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+city+"%")
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

	filters := gin.H{}
	if country != "" {
		filters["country"] = country
	}
	if supplierID != "" {
		filters["supplier_id"] = supplierID
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
	if city != "" {
		filters["city"] = city
	}
	if phone != "" {
		filters["phone"] = phone
	}
	if fax != "" {
		filters["fax"] = fax
	}

	c.JSON(http.StatusOK, gin.H{
		"filters": filters,
		"count":   len(suppliers),
		"data":    suppliers,
	})
}
