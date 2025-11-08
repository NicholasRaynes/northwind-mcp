package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-api/internal/db"
	"github.com/nicholasraynes/northwind-api/internal/models"
)

// GET /analytics/employee-performance?limit=10&year=1997
func GetEmployeePerformance(c *gin.Context) {
	limit := 10
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = val
		}
	}
	year := c.Query("year")

	query := `
		SELECT
			e.employee_id,
			(e.first_name || ' ' || e.last_name) AS full_name,
			e.title,
			e.country,
			COUNT(DISTINCT o.order_id) AS order_count,
			SUM(od.unit_price * od.quantity * (1 - od.discount)) AS total_revenue,
			AVG(od.unit_price * od.quantity * (1 - od.discount)) AS avg_order
		FROM order_details od
		JOIN orders o ON od.order_id = o.order_id
		JOIN employees e ON o.employee_id = e.employee_id
	`

	args := []interface{}{}
	if year != "" {
		query += " WHERE EXTRACT(YEAR FROM o.order_date) = $" + strconv.Itoa(len(args)+1)
		args = append(args, year)
	}

	query += `
		GROUP BY e.employee_id, e.first_name, e.last_name, e.title, e.country
		ORDER BY total_revenue DESC
		LIMIT ` + strconv.Itoa(limit)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	results := []models.EmployeePerformance{}
	for rows.Next() {
		var emp models.EmployeePerformance
		if err := rows.Scan(
			&emp.EmployeeID,
			&emp.FullName,
			&emp.Title,
			&emp.Country,
			&emp.OrderCount,
			&emp.TotalRevenue,
			&emp.AvgOrder,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		results = append(results, emp)
	}

	c.JSON(http.StatusOK, gin.H{
		"year":  year,
		"count": len(results),
		"data":  results,
	})
}
