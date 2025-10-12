package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-mcp/internal/db"
	"github.com/nicholasraynes/northwind-mcp/internal/models"
)

// GET /analytics/inventory-status?below_reorder=true&supplier=Exotic%20Liquids
func GetInventoryStatus(c *gin.Context) {
	belowReorder := c.DefaultQuery("below_reorder", "false") == "true"
	supplier := c.Query("supplier")
	category := c.Query("category")
	discontinued := c.DefaultQuery("discontinued", "")

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
	args := []interface{}{}
	where := []string{}

	if belowReorder {
		where = append(where, "p.units_in_stock <= p.reorder_level")
	}
	if supplier != "" {
		where = append(where, "LOWER(s.company_name) LIKE LOWER($"+strconv.Itoa(len(args)+1)+")")
		args = append(args, "%"+supplier+"%")
	}
	if category != "" {
		where = append(where, "LOWER(ca.category_name) LIKE LOWER($"+strconv.Itoa(len(args)+1)+")")
		args = append(args, "%"+category+"%")
	}
	if discontinued != "" {
		where = append(where, "p.discontinued = $"+strconv.Itoa(len(args)+1))
		args = append(args, discontinued == "true")
	}

	if len(where) > 0 {
		query += " WHERE " + joinConditions(where)
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

	c.JSON(http.StatusOK, gin.H{
		"below_reorder": belowReorder,
		"supplier":      supplier,
		"category":      category,
		"count":         len(results),
		"data":          results,
	})
}
