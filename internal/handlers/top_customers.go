package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-api/internal/db"
	"github.com/nicholasraynes/northwind-api/internal/models"
)

// GET /analytics/top-customers?limit=10&country=USA&year=1997
func GetTopCustomers(c *gin.Context) {
	limit := 10
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = val
		}
	}
	country := c.Query("country")
	year := c.Query("year")

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
	where := []string{}

	if country != "" {
		where = append(where, "c.country = $"+strconv.Itoa(len(args)+1))
		args = append(args, country)
	}
	if year != "" {
		where = append(where, "EXTRACT(YEAR FROM o.order_date) = $"+strconv.Itoa(len(args)+1))
		args = append(args, year)
	}

	if len(where) > 0 {
		query += " WHERE " + joinConditions(where)
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

	results := []models.TopCustomer{}
	for rows.Next() {
		var tc models.TopCustomer
		if err := rows.Scan(&tc.CustomerID, &tc.CompanyName, &tc.Country, &tc.TotalSales, &tc.OrderCount, &tc.AverageOrder); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		results = append(results, tc)
	}

	c.JSON(http.StatusOK, gin.H{
		"year":    year,
		"country": country,
		"count":   len(results),
		"data":    results,
	})
}

// helper function to safely join WHERE conditions
func joinConditions(conds []string) string {
	query := conds[0]
	for i := 1; i < len(conds); i++ {
		query += " AND " + conds[i]
	}
	return query
}
