package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-api/internal/db"
	"github.com/nicholasraynes/northwind-api/internal/models"
)

// GET /products
// Optional filters: ?category=Beverages&supplier=Exotic
func GetProducts(c *gin.Context) {
	category := c.Query("category")
	supplier := c.Query("supplier")
	name := c.Query("name")

	query := `
		SELECT
			p.product_id,
			p.product_name,
			p.supplier_id,
			s.company_name AS supplier_name,
			p.category_id,
			ca.category_name,
			p.quantity_per_unit,
			p.unit_price,
			p.units_in_stock,
			p.discontinued
		FROM products p
		JOIN suppliers s ON p.supplier_id = s.supplier_id
		JOIN categories ca ON p.category_id = ca.category_id
	`
	conditions := []string{}
	args := []interface{}{}

	if category != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(ca.category_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+category+"%")
	}
	if supplier != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(s.company_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+supplier+"%")
	}
	if name != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(p.product_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+name+"%")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY p.product_name"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		var p models.Product
		err := rows.Scan(
			&p.ProductID,
			&p.ProductName,
			&p.SupplierID,
			&p.SupplierName,
			&p.CategoryID,
			&p.CategoryName,
			&p.QuantityPerUnit,
			&p.UnitPrice,
			&p.UnitsInStock,
			&p.Discontinued,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		products = append(products, p)
	}

	c.JSON(http.StatusOK, gin.H{"count": len(products), "data": products})
}
