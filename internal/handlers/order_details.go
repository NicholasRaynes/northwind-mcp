package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-api/internal/db"
	"github.com/nicholasraynes/northwind-api/internal/models"
)

// GET /orders/details
// Optional filters: order_id, customer_id, product_id, product_name, category_name, supplier_name
func GetOrderDetails(c *gin.Context) {
	orderID := c.Query("order_id")
	customerID := c.Query("customer_id")
	productID := c.Query("product_id")
	productName := c.Query("product_name")
	categoryName := c.Query("category_name")
	supplierName := c.Query("supplier_name")

	query := `
		SELECT
			od.order_id,
			p.product_id,
			p.product_name,
			ca.category_name,
			s.company_name AS supplier_name,
			od.unit_price,
			od.quantity,
			od.discount,
			(od.unit_price * od.quantity * (1 - od.discount)) AS extended_price
		FROM order_details od
		JOIN products p ON od.product_id = p.product_id
		JOIN categories ca ON p.category_id = ca.category_id
		JOIN suppliers s ON p.supplier_id = s.supplier_id
		JOIN orders o ON od.order_id = o.order_id
	`

	conditions := []string{}
	args := []any{}

	if orderID != "" {
		conditions = append(conditions, fmt.Sprintf("od.order_id = $%d", len(args)+1))
		args = append(args, orderID)
	}
	if customerID != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(o.customer_id) = LOWER($%d)", len(args)+1))
		args = append(args, customerID)
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

	query += " ORDER BY od.order_id, p.product_name"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	results := []models.OrderDetail{}
	for rows.Next() {
		var d models.OrderDetail
		err := rows.Scan(
			&d.OrderID,
			&d.ProductID,
			&d.ProductName,
			&d.Category,
			&d.Supplier,
			&d.UnitPrice,
			&d.Quantity,
			&d.Discount,
			&d.ExtendedPrice,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		results = append(results, d)
	}

	filters := gin.H{}
	if orderID != "" {
		filters["order_id"] = orderID
	}
	if customerID != "" {
		filters["customer_id"] = customerID
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
