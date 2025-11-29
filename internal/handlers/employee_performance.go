package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-api/internal/db"
	"github.com/nicholasraynes/northwind-api/internal/models"
)

// GET /analytics/employee-performance
// Optional parameters: year, employee_id, full_name, title, country
func GetEmployeePerformance(c *gin.Context) {
	year := c.Query("year")
	employeeID := c.Query("employee_id")
	fullName := c.Query("full_name")
	job_title := c.Query("title")
	country := c.Query("country")

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

	args := []any{}
	conditions := []string{}

	if year != "" {
		conditions = append(conditions, fmt.Sprintf("EXTRACT(YEAR FROM o.order_date)::TEXT = $%d", len(args)+1))
		args = append(args, year)
	}
	if employeeID != "" {
		conditions = append(conditions, fmt.Sprintf("CAST(e.employee_id AS TEXT) LIKE $%d", len(args)+1))
		args = append(args, "%"+employeeID+"%")
	}
	if fullName != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(e.first_name || ' ' || e.last_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+fullName+"%")
	}
	if job_title != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(e.title) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+job_title+"%")
	}
	if country != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(e.country) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+country+"%")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += `
		GROUP BY e.employee_id, e.first_name, e.last_name, e.title, e.country
		ORDER BY total_revenue DESC
	`

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

	filters := gin.H{}
	if year != "" {
		filters["year"] = year
	}
	if employeeID != "" {
		filters["employee_id"] = employeeID
	}
	if fullName != "" {
		filters["full_name"] = fullName
	}
	if job_title != "" {
		filters["title"] = job_title
	}
	if country != "" {
		filters["country"] = country
	}

	c.JSON(http.StatusOK, gin.H{
		"filters": filters,
		"count":   len(results),
		"data":    results,
	})
}
