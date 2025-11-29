package handlers

import (
	"fmt" // Added for consistency with other files
	"net/http"
	"strings" // Added for string joining

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-api/internal/db"
	"github.com/nicholasraynes/northwind-api/internal/models"
)

// GET /analytics/top-products
// Optional parameters: year, product_id, product_name, category_name, supplier_name
func GetTopProducts(c *gin.Context) {
	year := c.Query("year")
	productID := c.Query("product_id")
	productName := c.Query("product_name")
	categoryName := c.Query("category_name")
	supplierName := c.Query("supplier_name")

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

	args := []any{}
	conditions := []string{}

	if year != "" {
		conditions = append(conditions, fmt.Sprintf("EXTRACT(YEAR FROM o.order_date)::TEXT = $%d", len(args)+1))
		args = append(args, year)
	}
	if productID != "" {
		conditions = append(conditions, fmt.Sprintf("CAST(p.product_id AS TEXT) LIKE $%d", len(args)+1))
		args = append(args, "%"+productID+"%")
	}
	if productName != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(p.product_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+productName+"%")
	}
	if categoryName != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(ca.category_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+categoryName+"%")
	}
	if supplierName != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(s.company_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+supplierName+"%")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += `
		GROUP BY p.product_id, p.product_name, ca.category_name, s.company_name
		ORDER BY total_revenue DESC
	`

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

	filters := gin.H{}
	if year != "" {
		filters["year"] = year
	}
	if productID != "" {
		filters["product_id"] = productID
	}
	if productName != "" {
		filters["product_name"] = productName
	}
	if categoryName != "" {
		filters["category_name"] = categoryName
	}
	if supplierName != "" {
		filters["supplier_name"] = supplierName
	}

	c.JSON(http.StatusOK, gin.H{
		"filters": filters,
		"count":   len(results),
		"data":    results,
	})
}
