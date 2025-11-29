package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-api/internal/db"
	"github.com/nicholasraynes/northwind-api/internal/models"
)

// GET /analytics/customer-retention?year=1997
// Optional parameters: year, customer_id, company_name, country, repeat_customer
func GetCustomerRetention(c *gin.Context) {
	year := c.Query("year")
	customerID := c.Query("customer_id")
	companyName := c.Query("company_name")
	country := c.Query("country")
	repeatCustomer := c.Query("repeat_customer")

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

	conditions := []string{}
	args := []any{}

	if year != "" {
		conditions = append(conditions, fmt.Sprintf("first_order_year <= $%d AND last_order_year >= $%d", len(args)+1, len(args)+2))
		args = append(args, year, year)
	}
	if customerID != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(customer_id) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+customerID+"%")
	}
	if companyName != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(company_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+companyName+"%")
	}
	if country != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(country) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+country+"%")
	}
	if repeatCustomer != "" {
		conditions = append(conditions, fmt.Sprintf("CAST(CASE WHEN active_years > 1 THEN true ELSE false END AS TEXT) LIKE $%d", len(args)+1))
		args = append(args, "%"+repeatCustomer+"%")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY repeat_customer DESC, active_years DESC, order_count DESC"

	rows, err := db.DB.Query(query, args...)
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

	filters := gin.H{}
	if year != "" {
		filters["year"] = year
	}
	if customerID != "" {
		filters["customer_id"] = customerID
	}
	if companyName != "" {
		filters["company_name"] = companyName
	}
	if country != "" {
		filters["country"] = country
	}
	if repeatCustomer != "" {
		filters["repeat_customer"] = repeatCustomer
	}
	filters["total_customers"] = totalCustomers
	filters["repeat_customers"] = repeatCount
	filters["retention_rate"] = retentionRate

	c.JSON(http.StatusOK, gin.H{
		"filters": filters,
		"count":   len(results),
		"data":    results,
	})
}
