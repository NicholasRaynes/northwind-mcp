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
// Optional parameters: customer_id, employee, year, country, order_id, customer_name, employee_name, order_date, required_date, shipped_date, ship_via, shipper_name, ship_name, ship_address, ship_city, ship_region, ship_postal_code, ship_country
func GetOrders(c *gin.Context) {
	customerID := c.Query("customer_id")
	employee := c.Query("employee")
	year := c.Query("year")
	country := c.Query("country")
	orderID := c.Query("order_id")
	customerName := c.Query("customer_name")
	employeeName := c.Query("employee_name")
	orderDate := c.Query("order_date")
	requiredDate := c.Query("required_date")
	shippedDate := c.Query("shipped_date")
	shipVia := c.Query("ship_via")
	shipperName := c.Query("shipper_name")
	shipName := c.Query("ship_name")
	shipAddress := c.Query("ship_address")
	shipCity := c.Query("ship_city")
	shipRegion := c.Query("ship_region")
	shipPostal := c.Query("ship_postal_code")
	shipCountry := c.Query("ship_country")

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
	args := []any{}

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
	if orderID != "" {
		conditions = append(conditions, fmt.Sprintf("CAST(o.order_id AS TEXT) LIKE $%d", len(args)+1))
		args = append(args, "%"+orderID+"%")
	}
	if customerName != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(c.company_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+customerName+"%")
	}
	if employeeName != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(e.first_name || ' ' || e.last_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+employeeName+"%")
	}
	if orderDate != "" {
		conditions = append(conditions, fmt.Sprintf("CAST(o.order_date AS TEXT) LIKE $%d", len(args)+1))
		args = append(args, "%"+orderDate+"%")
	}
	if requiredDate != "" {
		conditions = append(conditions, fmt.Sprintf("CAST(o.required_date AS TEXT) LIKE $%d", len(args)+1))
		args = append(args, "%"+requiredDate+"%")
	}
	if shippedDate != "" {
		conditions = append(conditions, fmt.Sprintf("CAST(o.shipped_date AS TEXT) LIKE $%d", len(args)+1))
		args = append(args, "%"+shippedDate+"%")
	}
	if shipVia != "" {
		conditions = append(conditions, fmt.Sprintf("CAST(o.ship_via AS TEXT) LIKE $%d", len(args)+1))
		args = append(args, "%"+shipVia+"%")
	}
	if shipperName != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(s.company_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+shipperName+"%")
	}
	if shipName != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(o.ship_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+shipName+"%")
	}
	if shipAddress != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(o.ship_address) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+shipAddress+"%")
	}
	if shipCity != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(o.ship_city) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+shipCity+"%")
	}
	if shipRegion != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(o.ship_region) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+shipRegion+"%")
	}
	if shipPostal != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(o.ship_postal_code) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+shipPostal+"%")
	}
	if shipCountry != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(o.ship_country) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+shipCountry+"%")
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

	filters := gin.H{}
	if customerID != "" {
		filters["customer_id"] = customerID
	}
	if employee != "" {
		filters["employee"] = employee
	}
	if year != "" {
		filters["year"] = year
	}
	if country != "" {
		filters["country"] = country
	}
	if orderID != "" {
		filters["order_id"] = orderID
	}
	if customerName != "" {
		filters["customer_name"] = customerName
	}
	if employeeName != "" {
		filters["employee_name"] = employeeName
	}
	if orderDate != "" {
		filters["order_date"] = orderDate
	}
	if requiredDate != "" {
		filters["required_date"] = requiredDate
	}
	if shippedDate != "" {
		filters["shipped_date"] = shippedDate
	}
	if shipVia != "" {
		filters["ship_via"] = shipVia
	}
	if shipperName != "" {
		filters["shipper_name"] = shipperName
	}
	if shipName != "" {
		filters["ship_name"] = shipName
	}
	if shipAddress != "" {
		filters["ship_address"] = shipAddress
	}
	if shipCity != "" {
		filters["ship_city"] = shipCity
	}
	if shipRegion != "" {
		filters["ship_region"] = shipRegion
	}
	if shipPostal != "" {
		filters["ship_postal_code"] = shipPostal
	}
	if shipCountry != "" {
		filters["ship_country"] = shipCountry
	}

	c.JSON(http.StatusOK, gin.H{
		"filters": filters,
		"count":   len(orders),
		"data":    orders,
	})
}
