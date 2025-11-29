package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-api/internal/db"
	"github.com/nicholasraynes/northwind-api/internal/models"
)

// GET /analytics/customer-orders
// Optional parameters: customer_id, year, order_id, company_name, order_date, shipped_date, country
func GetCustomerOrders(c *gin.Context) {
	customerID := c.Query("customer_id")
	year := c.Query("year")
	orderID := c.Query("order_id")
	companyName := c.Query("company_name")
	orderDate := c.Query("order_date")
	shippedDate := c.Query("shipped_date")
	country := c.Query("country")

	query := `
		SELECT
			c.customer_id,
			c.company_name,
			o.order_id,
			TO_CHAR(o.order_date, 'YYYY-MM-DD') AS order_date,
			SUM(od.unit_price * od.quantity * (1 - od.discount)) AS total_amount,
			COALESCE(TO_CHAR(o.shipped_date, 'YYYY-MM-DD'), '') AS shipped_date,
			c.country
		FROM order_details od
		JOIN orders o ON od.order_id = o.order_id
		JOIN customers c ON o.customer_id = c.customer_id
	`

	args := []any{}
	conditions := []string{}

	if customerID != "" {
		conditions = append(conditions, fmt.Sprintf("c.customer_id = $%d", len(args)+1))
		args = append(args, customerID)
	}
	if year != "" {
		conditions = append(conditions, fmt.Sprintf("EXTRACT(YEAR FROM o.order_date)::TEXT = $%d", len(args)+1))
		args = append(args, year)
	}
	if orderID != "" {
		conditions = append(conditions, fmt.Sprintf("CAST(o.order_id AS TEXT) LIKE $%d", len(args)+1))
		args = append(args, "%"+orderID+"%")
	}
	if companyName != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(c.company_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+companyName+"%")
	}
	if orderDate != "" {
		conditions = append(conditions, fmt.Sprintf("CAST(TO_CHAR(o.order_date, 'YYYY-MM-DD') AS TEXT) LIKE $%d", len(args)+1))
		args = append(args, "%"+orderDate+"%")
	}
	if shippedDate != "" {
		conditions = append(conditions, fmt.Sprintf("CAST(COALESCE(TO_CHAR(o.shipped_date, 'YYYY-MM-DD'), '') AS TEXT) LIKE $%d", len(args)+1))
		args = append(args, "%"+shippedDate+"%")
	}
	if country != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(c.country) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+country+"%")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += `
		GROUP BY c.customer_id, c.company_name, o.order_id, o.order_date, o.shipped_date, c.country
		ORDER BY o.order_date DESC
	`

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	results := []models.CustomerOrderSummary{}
	for rows.Next() {
		var co models.CustomerOrderSummary
		if err := rows.Scan(
			&co.CustomerID,
			&co.CompanyName,
			&co.OrderID,
			&co.OrderDate,
			&co.TotalAmount,
			&co.ShippedDate,
			&co.Country,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		results = append(results, co)
	}

	filters := gin.H{}
	if customerID != "" {
		filters["customer_id"] = customerID
	}
	if year != "" {
		filters["year"] = year
	}
	if orderID != "" {
		filters["order_id"] = orderID
	}
	if companyName != "" {
		filters["company_name"] = companyName
	}
	if orderDate != "" {
		filters["order_date"] = orderDate
	}
	if shippedDate != "" {
		filters["shipped_date"] = shippedDate
	}
	if country != "" {
		filters["country"] = country
	}

	c.JSON(http.StatusOK, gin.H{
		"filters": filters,
		"count":   len(results),
		"data":    results,
	})
}
