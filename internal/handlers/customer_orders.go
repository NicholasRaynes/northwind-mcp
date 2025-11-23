package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-api/internal/db"
	"github.com/nicholasraynes/northwind-api/internal/models"
)

// GET /analytics/customer-orders?customer_id=ALFKI&year=1997
func GetCustomerOrders(c *gin.Context) {
	customerID := c.Query("customer_id")
	year := c.Query("year")

	query := `
		SELECT
			c.customer_id,
			c.company_name,
			o.order_id,
			TO_CHAR(o.order_date, 'YYYY-MM-DD') AS order_date,
			SUM(od.unit_price * od.quantity * (1 - od.discount)) AS total_amount,
			TO_CHAR(o.shipped_date, 'YYYY-MM-DD') AS shipped_date,
			c.country
		FROM order_details od
		JOIN orders o ON od.order_id = o.order_id
		JOIN customers c ON o.customer_id = c.customer_id
	`
	args := []any{}
	where := []string{}

	if customerID != "" {
		where = append(where, "c.customer_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, customerID)
	}
	if year != "" {
		where = append(where, "EXTRACT(YEAR FROM o.order_date) = $"+strconv.Itoa(len(args)+1))
		args = append(args, year)
	}

	if len(where) > 0 {
		query += " WHERE " + joinConditions(where)
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

	c.JSON(http.StatusOK, gin.H{
		"customer_id": customerID,
		"year":        year,
		"count":       len(results),
		"data":        results,
	})
}
