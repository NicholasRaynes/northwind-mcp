package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-api/internal/db"
	"github.com/nicholasraynes/northwind-api/internal/models"
)

// GET /orders
// Optional filters: ?customer_id=ALFKI&employee=Fuller&year=1997
func GetOrders(c *gin.Context) {
	customerID := c.Query("customer_id")
	employee := c.Query("employee")
	year := c.Query("year")
	country := c.Query("country")

	query := `
		SELECT
			o.order_id,
			o.customer_id,
			c.company_name AS customer_name,
			(e.first_name || ' ' || e.last_name) AS employee_name,
			o.order_date,
			o.required_date,
			o.shipped_date,
			o.ship_via,
			s.company_name AS shipper_name,
			o.freight,
			o.ship_name,
			o.ship_address,
			o.ship_city,
			o.ship_region,
			o.ship_postal_code,
			o.ship_country
		FROM orders o
		JOIN customers c ON o.customer_id = c.customer_id
		JOIN employees e ON o.employee_id = e.employee_id
		JOIN shippers s ON o.ship_via = s.shipper_id
	`

	conditions := []string{}
	args := []interface{}{}

	if customerID != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(o.customer_id) = LOWER($%d)", len(args)+1))
		args = append(args, customerID)
	}

	if employee != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(e.first_name || ' ' || e.last_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+employee+"%")
	}

	if year != "" {
		conditions = append(conditions, fmt.Sprintf("EXTRACT(YEAR FROM o.order_date) = $%d", len(args)+1))
		args = append(args, year)
	}

	if country != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(o.ship_country) = LOWER($%d)", len(args)+1))
		args = append(args, country)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY o.order_date DESC"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	orders := []models.Order{}

	for rows.Next() {
		var o models.Order
		err := rows.Scan(
			&o.OrderID,
			&o.CustomerID,
			&o.CustomerName,
			&o.EmployeeName,
			&o.OrderDate,
			&o.RequiredDate,
			&o.ShippedDate,
			&o.ShipVia,
			&o.ShipperName,
			&o.Freight,
			&o.ShipName,
			&o.ShipAddress,
			&o.ShipCity,
			&o.ShipRegion,
			&o.ShipPostal,
			&o.ShipCountry,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		orders = append(orders, o)
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(orders),
		"data":  orders,
	})
}
