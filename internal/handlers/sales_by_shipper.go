package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-mcp/internal/db"
	"github.com/nicholasraynes/northwind-mcp/internal/models"
)

// GET /summary/sales-by-shipper?year=1997
func GetSalesByShipper(c *gin.Context) {
	year := c.Query("year")

	query := `
		SELECT
			s.company_name AS group_key,
			SUM(od.unit_price * od.quantity * (1 - od.discount)) AS total_sales,
			COUNT(DISTINCT o.order_id) AS order_count
		FROM order_details od
		JOIN orders o ON od.order_id = o.order_id
		JOIN shippers s ON o.ship_via = s.shipper_id
	`
	args := []interface{}{}
	if year != "" {
		query += " WHERE EXTRACT(YEAR FROM o.order_date) = $1"
		args = append(args, year)
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

	c.JSON(http.StatusOK, gin.H{
		"year":  year,
		"count": len(results),
		"data":  results,
	})
}
