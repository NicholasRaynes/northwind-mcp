package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-api/internal/db"
	"github.com/nicholasraynes/northwind-api/internal/models"
)

// GET /summary/sales-by-employee
// Optional parameters: year, employee_name
func GetSalesByEmployee(c *gin.Context) {
	year := c.Query("year")
	employeeName := c.Query("employee_name")

	query := `
		SELECT
			(e.first_name || ' ' || e.last_name) AS employee_name,
			SUM(od.unit_price * od.quantity * (1 - od.discount)) AS total_sales,
			COUNT(DISTINCT o.order_id) AS order_count
		FROM order_details od
		JOIN orders o ON od.order_id = o.order_id
		JOIN employees e ON o.employee_id = e.employee_id
	`

	conditions := []string{}
	args := []any{}

	if year != "" {
		conditions = append(conditions, fmt.Sprintf("EXTRACT(YEAR FROM o.order_date)::TEXT = $%d", len(args)+1))
		args = append(args, year)
	}
	if employeeName != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(e.first_name || ' ' || e.last_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+employeeName+"%")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += `
		GROUP BY e.first_name, e.last_name
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
	if employeeName != "" {
		filters["employee_name"] = employeeName
	}

	c.JSON(http.StatusOK, gin.H{
		"filters": filters,
		"count":   len(results),
		"data":    results,
	})
}
