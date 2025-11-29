package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-api/internal/db"
	"github.com/nicholasraynes/northwind-api/internal/models"
)

// GET /analytics/supplier-performance
// Optional parameters: year, supplier_id, supplier_name, country, top_category
func GetSupplierPerformance(c *gin.Context) {
	year := c.Query("year")
	supplierID := c.Query("supplier_id")
	supplierName := c.Query("supplier_name")
	country := c.Query("country")
	topCategory := c.Query("top_category")

	args := []any{}
	cteConditions := []string{}

	if year != "" {
		cteConditions = append(cteConditions, fmt.Sprintf("EXTRACT(YEAR FROM o.order_date)::TEXT = $%d", len(args)+1))
		args = append(args, year)
	}

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

	if len(cteConditions) > 0 {
		query += " WHERE " + strings.Join(cteConditions, " AND ")
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
	`

	finalConditions := []string{}

	if supplierID != "" {
		finalConditions = append(finalConditions, fmt.Sprintf("CAST(supplier_id AS TEXT) LIKE $%d", len(args)+1))
		args = append(args, "%"+supplierID+"%")
	}
	if supplierName != "" {
		finalConditions = append(finalConditions, fmt.Sprintf("LOWER(supplier_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+supplierName+"%")
	}
	if country != "" {
		finalConditions = append(finalConditions, fmt.Sprintf("LOWER(country) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+country+"%")
	}
	if topCategory != "" {
		finalConditions = append(finalConditions, fmt.Sprintf("LOWER(category_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+topCategory+"%")
	}

	if len(finalConditions) > 0 {
		query += " AND " + strings.Join(finalConditions, " AND ")
	}

	query += " ORDER BY total_revenue DESC"

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

	filters := gin.H{}
	if year != "" {
		filters["year"] = year
	}
	if supplierID != "" {
		filters["supplier_id"] = supplierID
	}
	if supplierName != "" {
		filters["supplier_name"] = supplierName
	}
	if country != "" {
		filters["country"] = country
	}
	if topCategory != "" {
		filters["top_category"] = topCategory
	}

	c.JSON(http.StatusOK, gin.H{
		"filters": filters,
		"count":   len(results),
		"data":    results,
	})
}
