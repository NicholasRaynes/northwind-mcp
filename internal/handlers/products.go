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
// Optional parameters: product_id, product_name, supplier_id, supplier_name, category_id, category_name, discontinued
func GetProducts(c *gin.Context) {
	productID := c.Query("product_id")
	productName := c.Query("product_name")
	supplierID := c.Query("supplier_id")
	supplierName := c.Query("supplier_name")
	categoryID := c.Query("category_id")
	categoryName := c.Query("category_name")
	discontinued := c.Query("discontinued")

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
	args := []any{}

	if productID != "" {
		conditions = append(conditions, fmt.Sprintf("CAST(p.product_id AS TEXT) LIKE $%d", len(args)+1))
		args = append(args, "%"+productID+"%")
	}
	if productName != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(p.product_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+productName+"%")
	}
	if supplierID != "" {
		conditions = append(conditions, fmt.Sprintf("CAST(p.supplier_id AS TEXT) LIKE $%d", len(args)+1))
		args = append(args, "%"+supplierID+"%")
	}
	if supplierName != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(s.company_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+supplierName+"%")
	}
	if categoryID != "" {
		conditions = append(conditions, fmt.Sprintf("CAST(p.category_id AS TEXT) LIKE $%d", len(args)+1))
		args = append(args, "%"+categoryID+"%")
	}
	if categoryName != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(ca.category_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+categoryName+"%")
	}
	if discontinued != "" {
		conditions = append(conditions, fmt.Sprintf("CAST(p.discontinued AS TEXT) LIKE $%d", len(args)+1))
		args = append(args, "%"+discontinued+"%")
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

	filters := gin.H{}
	if productID != "" {
		filters["product_id"] = productID
	}
	if productName != "" {
		filters["product_name"] = productName
	}
	if supplierID != "" {
		filters["supplier_id"] = supplierID
	}
	if supplierName != "" {
		filters["supplier_name"] = supplierName
	}
	if categoryID != "" {
		filters["category_id"] = categoryID
	}
	if categoryName != "" {
		filters["category_name"] = categoryName
	}
	if discontinued != "" {
		filters["discontinued"] = discontinued
	}

	c.JSON(http.StatusOK, gin.H{
		"filters": filters,
		"count":   len(products),
		"data":    products,
	})
}
