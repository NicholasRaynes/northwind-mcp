package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-api/internal/db"
	"github.com/nicholasraynes/northwind-api/internal/models"
)

// GET /analytics/supplier-performance?limit=10&year=1997
func GetSupplierPerformance(c *gin.Context) {
	limit := 10
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = val
		}
	}
	year := c.Query("year")

	query := `
		WITH supplier_stats AS (
			SELECT
				s.supplier_id,
				s.company_name AS supplier_name,
				s.country,
				COUNT(DISTINCT p.product_id) AS product_count,
				SUM(od.quantity) AS units_sold,
				SUM(od.unit_price * od.quantity * (1 - od.discount)) AS total_revenue,
				AVG(od.unit_price) AS average_price,
				ca.category_name,
				ROW_NUMBER() OVER (PARTITION BY s.supplier_id ORDER BY SUM(od.unit_price * od.quantity) DESC) AS cat_rank
			FROM order_details od
			JOIN orders o ON od.order_id = o.order_id
			JOIN products p ON od.product_id = p.product_id
			JOIN categories ca ON p.category_id = ca.category_id
			JOIN suppliers s ON p.supplier_id = s.supplier_id
	`
	args := []any{}
	if year != "" {
		query += " WHERE EXTRACT(YEAR FROM o.order_date) = $" + strconv.Itoa(len(args)+1)
		args = append(args, year)
	}
	query += `
			GROUP BY s.supplier_id, s.company_name, s.country, ca.category_name
		)
		SELECT
			supplier_id,
			supplier_name,
			country,
			product_count,
			units_sold,
			total_revenue,
			average_price,
			category_name AS top_category
		FROM supplier_stats
		WHERE cat_rank = 1
		ORDER BY total_revenue DESC
		LIMIT ` + strconv.Itoa(limit)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	results := []models.SupplierPerformance{}
	for rows.Next() {
		var sp models.SupplierPerformance
		if err := rows.Scan(
			&sp.SupplierID,
			&sp.SupplierName,
			&sp.Country,
			&sp.ProductCount,
			&sp.UnitsSold,
			&sp.TotalRevenue,
			&sp.AveragePrice,
			&sp.TopCategory,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		results = append(results, sp)
	}

	c.JSON(http.StatusOK, gin.H{
		"year":  year,
		"count": len(results),
		"data":  results,
	})
}
