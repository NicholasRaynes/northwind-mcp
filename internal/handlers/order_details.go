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
// Optional filters: ?order_id=10248&product=Chai&customer_id=ALFKI
func GetOrderDetails(c *gin.Context) {
	orderID := c.Query("order_id")
	product := c.Query("product")
	customerID := c.Query("customer_id")

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
	args := []interface{}{}

	if orderID != "" {
		conditions = append(conditions, fmt.Sprintf("od.order_id = $%d", len(args)+1))
		args = append(args, orderID)
	}
	if product != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(p.product_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+product+"%")
	}
	if customerID != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(o.customer_id) = LOWER($%d)", len(args)+1))
		args = append(args, customerID)
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

	c.JSON(http.StatusOK, gin.H{
		"count": len(results),
		"data":  results,
	})
}
