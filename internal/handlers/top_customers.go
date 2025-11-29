package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-api/internal/db"
	"github.com/nicholasraynes/northwind-api/internal/models"
)

// GET /analytics/top-customers
// Optional parameters: country, year, customer_id, company_name
func GetTopCustomers(c *gin.Context) {
	country := c.Query("country")
	year := c.Query("year")
	customerID := c.Query("customer_id")
	companyName := c.Query("company_name")

	query := `
		SELECT
			c.customer_id,
			c.company_name,
			c.country,
			SUM(od.unit_price * od.quantity * (1 - od.discount)) AS total_sales,
			COUNT(DISTINCT o.order_id) AS order_count,
			AVG(od.unit_price * od.quantity * (1 - od.discount)) AS average_order
		FROM order_details od
		JOIN orders o ON od.order_id = o.order_id
		JOIN customers c ON o.customer_id = c.customer_id
	`

	args := []any{}
	conditions := []string{}

	if country != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(c.country) LIKE LOWER($%d)", len(args)+1))
		args = append(args, country)
	}
	if year != "" {
		conditions = append(conditions, fmt.Sprintf("EXTRACT(YEAR FROM o.order_date)::TEXT = $%d", len(args)+1))
		args = append(args, year)
	}
	if customerID != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(c.customer_id) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+customerID+"%")
	}
	if companyName != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(c.company_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+companyName+"%")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += `
		GROUP BY c.customer_id, c.company_name, c.country
		ORDER BY total_sales DESC
	`

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	results := []models.TopCustomer{}
	for rows.Next() {
		var tc models.TopCustomer
		if err := rows.Scan(&tc.CustomerID, &tc.CompanyName, &tc.Country, &tc.TotalSales, &tc.OrderCount, &tc.AverageOrder); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		results = append(results, tc)
	}

	filters := gin.H{}
	if year != "" {
		filters["year"] = year
	}
	if country != "" {
		filters["country"] = country
	}
	if customerID != "" {
		filters["customer_id"] = customerID
	}
	if companyName != "" {
		filters["company_name"] = companyName
	}

	c.JSON(http.StatusOK, gin.H{
		"filters": filters,
		"count":   len(results),
		"data":    results,
	})
}
