package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-api/internal/db"
	"github.com/nicholasraynes/northwind-api/internal/models"
)

// GET /summary/sales-by-category
// Optional parameters: year, category_name
func GetSalesByCategory(c *gin.Context) {
	year := c.Query("year")
	categoryName := c.Query("category_name")

	query := `
		SELECT
			COALESCE(ca.category_name, 'Unknown') AS category_name,
			SUM(od.unit_price * od.quantity * (1 - od.discount)) AS total_sales,
			COUNT(DISTINCT o.order_id) AS order_count
		FROM order_details od
		JOIN orders o ON od.order_id = o.order_id
		JOIN products p ON od.product_id = p.product_id
		JOIN categories ca ON p.category_id = ca.category_id
	`

	conditions := []string{}
	args := []any{}

	if year != "" {
		conditions = append(conditions, fmt.Sprintf("EXTRACT(YEAR FROM o.order_date)::TEXT = $%d", len(args)+1))
		args = append(args, year)
	}
	if categoryName != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(COALESCE(ca.category_name, 'Unknown')) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+categoryName+"%")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += `
		GROUP BY ca.category_name
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
	if categoryName != "" {
		filters["category_name"] = categoryName
	}

	c.JSON(http.StatusOK, gin.H{
		"filters": filters,
		"count":   len(results),
		"data":    results,
	})
}
