package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-api/internal/db"
	"github.com/nicholasraynes/northwind-api/internal/models"
)

// GET /analytics/customer-retention?year=1997
func GetCustomerRetention(c *gin.Context) {
	year := c.Query("year")

	query := `
		WITH customer_years AS (
			SELECT
				c.customer_id,
				c.company_name,
				c.country,
				MIN(EXTRACT(YEAR FROM o.order_date))::int AS first_order_year,
				MAX(EXTRACT(YEAR FROM o.order_date))::int AS last_order_year,
				COUNT(DISTINCT o.order_id) AS order_count,
				COUNT(DISTINCT EXTRACT(YEAR FROM o.order_date)) AS active_years
			FROM orders o
			JOIN customers c ON o.customer_id = c.customer_id
			GROUP BY c.customer_id, c.company_name, c.country
		)
		SELECT
			customer_id,
			company_name,
			country,
			first_order_year,
			last_order_year,
			order_count,
			active_years,
			CASE WHEN active_years > 1 THEN true ELSE false END AS repeat_customer
		FROM customer_years
	`

	// Optional filter by year
	if year != "" {
		query += " WHERE first_order_year <= " + year + " AND last_order_year >= " + year
	}

	query += " ORDER BY repeat_customer DESC, active_years DESC, order_count DESC"

	rows, err := db.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	results := []models.CustomerRetention{}
	for rows.Next() {
		var cr models.CustomerRetention
		if err := rows.Scan(
			&cr.CustomerID,
			&cr.CompanyName,
			&cr.Country,
			&cr.FirstOrderYear,
			&cr.LastOrderYear,
			&cr.OrderCount,
			&cr.ActiveYears,
			&cr.RepeatCustomer,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		results = append(results, cr)
	}

	// Compute overall metrics
	totalCustomers := len(results)
	repeatCount := 0
	for _, r := range results {
		if r.RepeatCustomer {
			repeatCount++
		}
	}

	retentionRate := 0.0
	if totalCustomers > 0 {
		retentionRate = float64(repeatCount) / float64(totalCustomers)
	}

	c.JSON(http.StatusOK, gin.H{
		"year":             year,
		"total_customers":  totalCustomers,
		"repeat_customers": repeatCount,
		"retention_rate":   retentionRate,
		"data":             results,
	})
}
