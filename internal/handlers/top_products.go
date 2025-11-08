package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-api/internal/db"
	"github.com/nicholasraynes/northwind-api/internal/models"
)

// GET /analytics/top-products?limit=10&year=1997&category=Beverages
func GetTopProducts(c *gin.Context) {
	limit := 10
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = val
		}
	}
	year := c.Query("year")
	category := c.Query("category")

	query := `
		SELECT
			p.product_id,
			p.product_name,
			ca.category_name,
			s.company_name AS supplier_name,
			SUM(od.quantity) AS units_sold,
			SUM(od.unit_price * od.quantity * (1 - od.discount)) AS total_revenue,
			AVG(od.unit_price) AS average_price
		FROM order_details od
		JOIN orders o ON od.order_id = o.order_id
		JOIN products p ON od.product_id = p.product_id
		JOIN categories ca ON p.category_id = ca.category_id
		JOIN suppliers s ON p.supplier_id = s.supplier_id
	`

	args := []interface{}{}
	where := []string{}

	if year != "" {
		where = append(where, "EXTRACT(YEAR FROM o.order_date) = $"+strconv.Itoa(len(args)+1))
		args = append(args, year)
	}
	if category != "" {
		where = append(where, "LOWER(ca.category_name) = LOWER($"+strconv.Itoa(len(args)+1)+")")
		args = append(args, category)
	}

	if len(where) > 0 {
		query += " WHERE " + joinConditions(where)
	}

	query += `
		GROUP BY p.product_id, p.product_name, ca.category_name, s.company_name
		ORDER BY total_revenue DESC
		LIMIT ` + strconv.Itoa(limit)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	results := []models.TopProduct{}
	for rows.Next() {
		var tp models.TopProduct
		if err := rows.Scan(
			&tp.ProductID,
			&tp.ProductName,
			&tp.CategoryName,
			&tp.SupplierName,
			&tp.UnitsSold,
			&tp.TotalRevenue,
			&tp.AveragePrice,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		results = append(results, tp)
	}

	c.JSON(http.StatusOK, gin.H{
		"year":     year,
		"category": category,
		"count":    len(results),
		"data":     results,
	})
}
