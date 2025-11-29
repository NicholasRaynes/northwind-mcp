package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-api/internal/db"
	"github.com/nicholasraynes/northwind-api/internal/models"
)

// GET /summary/sales-by-shipper
// Optional parameters: year, company_name
func GetSalesByShipper(c *gin.Context) {
	year := c.Query("year")
	companyName := c.Query("company_name")

	query := `
		SELECT
			s.company_name AS company_name,
			SUM(od.unit_price * od.quantity * (1 - od.discount)) AS total_sales,
			COUNT(DISTINCT o.order_id) AS order_count
		FROM order_details od
		JOIN orders o ON od.order_id = o.order_id
		JOIN shippers s ON o.ship_via = s.shipper_id
	`

	conditions := []string{}
	args := []any{}

	if year != "" {
		conditions = append(conditions, fmt.Sprintf("EXTRACT(YEAR FROM o.order_date)::TEXT = $%d", len(args)+1))
		args = append(args, year)
	}
	if companyName != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(s.company_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+companyName+"%")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += `
		GROUP BY s.company_name
		ORDER BY total_sales DESC
	`

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	results := []models.SalesSummary{}
	for rows.Next() {
		var s models.SalesSummary
		if err := rows.Scan(&s.GroupKey, &s.TotalSales, &s.OrderCount); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		results = append(results, s)
	}

	filters := gin.H{}
	if year != "" {
		filters["year"] = year
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
