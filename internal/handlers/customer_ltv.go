package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-mcp/internal/db"
	"github.com/nicholasraynes/northwind-mcp/internal/models"
)

// GET /analytics/customer-ltv?limit=10&country=Germany
func GetCustomerLTV(c *gin.Context) {
	limit := 10
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = val
		}
	}
	country := c.Query("country")

	query := `
		SELECT
			c.customer_id,
			c.company_name,
			c.country,
			MIN(TO_CHAR(o.order_date, 'YYYY-MM-DD')) AS first_order,
			MAX(TO_CHAR(o.order_date, 'YYYY-MM-DD')) AS last_order,
			COUNT(DISTINCT o.order_id) AS order_count,
			SUM(od.unit_price * od.quantity * (1 - od.discount)) AS total_sales,
			AVG(od.unit_price * od.quantity * (1 - od.discount)) AS avg_order
		FROM order_details od
		JOIN orders o ON od.order_id = o.order_id
		JOIN customers c ON o.customer_id = c.customer_id
	`
	args := []interface{}{}
	if country != "" {
		query += " WHERE c.country = $1"
		args = append(args, country)
	}
	query += `
		GROUP BY c.customer_id, c.company_name, c.country
		ORDER BY total_sales DESC
		LIMIT ` + strconv.Itoa(limit)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	results := []models.CustomerLTV{}
	for rows.Next() {
		var cl models.CustomerLTV
		if err := rows.Scan(
			&cl.CustomerID,
			&cl.CompanyName,
			&cl.Country,
			&cl.FirstOrder,
			&cl.LastOrder,
			&cl.OrderCount,
			&cl.TotalSales,
			&cl.AvgOrder,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		results = append(results, cl)
	}

	c.JSON(http.StatusOK, gin.H{
		"country": country,
		"count":   len(results),
		"data":    results,
	})
}
