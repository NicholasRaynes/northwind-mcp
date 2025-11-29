package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-api/internal/db"
	"github.com/nicholasraynes/northwind-api/internal/models"
)

// GET /analytics/inventory-status
// Optional parameters: product_id, product_name, supplier_name, category_name, discontinued, needs_reorder
func GetInventoryStatus(c *gin.Context) {
	productID := c.Query("product_id")
	productName := c.Query("product_name")
	supplierName := c.Query("supplier_name")
	categoryName := c.Query("category_name")
	discontinued := c.Query("discontinued")
	needsReorder := c.Query("needs_reorder")

	query := `
		SELECT
			p.product_id,
			p.product_name,
			s.company_name AS supplier_name,
			ca.category_name,
			p.units_in_stock,
			p.reorder_level,
			p.discontinued,
			CASE WHEN p.units_in_stock <= p.reorder_level THEN true ELSE false END AS needs_reorder
		FROM products p
		JOIN suppliers s ON p.supplier_id = s.supplier_id
		JOIN categories ca ON p.category_id = ca.category_id
	`

	args := []any{}
	conditions := []string{}

	if productID != "" {
		conditions = append(conditions, fmt.Sprintf("CAST(p.product_id AS TEXT) LIKE $%d", len(args)+1))
		args = append(args, "%"+productID+"%")
	}
	if productName != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(p.product_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+productName+"%")
	}
	if supplierName != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(s.company_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+supplierName+"%")
	}
	if categoryName != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(ca.category_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+categoryName+"%")
	}
	if discontinued != "" {
		conditions = append(conditions, fmt.Sprintf("CAST(p.discontinued AS TEXT) LIKE $%d", len(args)+1))
		args = append(args, "%"+discontinued+"%")
	}
	if needsReorder != "" {
		conditions = append(conditions, fmt.Sprintf("CAST(CASE WHEN p.units_in_stock <= p.reorder_level THEN true ELSE false END AS TEXT) LIKE $%d", len(args)+1))
		args = append(args, "%"+needsReorder+"%")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += `
		ORDER BY needs_reorder DESC, p.units_in_stock ASC
	`

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	results := []models.InventoryStatus{}
	for rows.Next() {
		var inv models.InventoryStatus
		if err := rows.Scan(
			&inv.ProductID,
			&inv.ProductName,
			&inv.SupplierName,
			&inv.CategoryName,
			&inv.UnitsInStock,
			&inv.ReorderLevel,
			&inv.Discontinued,
			&inv.NeedsReorder,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		results = append(results, inv)
	}

	filters := gin.H{}
	if productID != "" {
		filters["product_id"] = productID
	}
	if productName != "" {
		filters["product_name"] = productName
	}
	if supplierName != "" {
		filters["supplier_name"] = supplierName
	}
	if categoryName != "" {
		filters["category_name"] = categoryName
	}
	if discontinued != "" {
		filters["discontinued"] = discontinued
	}
	if needsReorder != "" {
		filters["needs_reorder"] = needsReorder
	}

	c.JSON(http.StatusOK, gin.H{
		"filters": filters,
		"count":   len(results),
		"data":    results,
	})
}
