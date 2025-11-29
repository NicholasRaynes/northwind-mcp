package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-api/internal/db"
	"github.com/nicholasraynes/northwind-api/internal/models"
)

// GET /summary/sales-by-year
func GetSalesByYear(c *gin.Context) {
	query := `
		SELECT
			EXTRACT(YEAR FROM o.order_date)::text AS year,
			SUM(od.unit_price * od.quantity * (1 - od.discount)) AS total_sales,
			COUNT(DISTINCT o.order_id) AS order_count
		FROM order_details od
		JOIN orders o ON od.order_id = o.order_id
	`

	query += `
		GROUP BY year
		ORDER BY year
	`

	rows, err := db.DB.Query(query)
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
		"filters": gin.H{},
		"count":   len(results),
		"data":    results,
	})
}
