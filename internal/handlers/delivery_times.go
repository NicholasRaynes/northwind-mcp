package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-api/internal/db"
	"github.com/nicholasraynes/northwind-api/internal/models"
)

// GET /analytics/delivery-times?year=1997&group=shipper
func GetDeliveryTimes(c *gin.Context) {
	group := c.DefaultQuery("group", "shipper") // shipper or employee
	year := c.Query("year")

	groupField := "s.company_name"
	joinClause := "JOIN shippers s ON o.ship_via = s.shipper_id"

	if group == "employee" {
		groupField = "(e.first_name || ' ' || e.last_name)"
		joinClause = "JOIN employees e ON o.employee_id = e.employee_id"
	}

	query := `
		SELECT
			` + groupField + ` AS name,
			COUNT(o.order_id) AS total_orders,
			AVG(EXTRACT(DAY FROM (o.shipped_date - o.order_date))) AS avg_delivery_days,
			MAX(EXTRACT(DAY FROM (o.shipped_date - o.order_date))) AS max_delivery_days,
			MIN(EXTRACT(DAY FROM (o.shipped_date - o.order_date))) AS min_delivery_days,
			SUM(CASE WHEN o.required_date IS NOT NULL AND o.shipped_date > o.required_date THEN 1 ELSE 0 END) AS late_shipments,
			(1 - SUM(CASE WHEN o.required_date IS NOT NULL AND o.shipped_date > o.required_date THEN 1 ELSE 0 END)::float / COUNT(o.order_id)) AS on_time_rate
		FROM orders o
		` + joinClause

	args := []any{}
	where := []string{}

	where = append(where, "o.shipped_date IS NOT NULL")

	if year != "" {
		yearInt, err := strconv.Atoi(year)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year parameter"})
			return
		}
		args = append(args, yearInt)
		where = append(where, "EXTRACT(YEAR FROM o.order_date)::int = $"+strconv.Itoa(len(args)))
	}

	if len(where) > 0 {
		query += " WHERE "

		for i, condition := range where {
			if i > 0 {
				query += " AND "
			}
			query += condition
		}
	}

	query += `
		GROUP BY ` + groupField + `
		ORDER BY avg_delivery_days ASC
	`

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	results := []models.DeliveryTimes{}
	for rows.Next() {
		var d models.DeliveryTimes
		if group == "shipper" {
			if err := rows.Scan(
				&d.ShipperName,
				&d.TotalOrders,
				&d.AvgDeliveryDays,
				&d.MaxDeliveryDays,
				&d.MinDeliveryDays,
				&d.LateShipments,
				&d.OnTimeRate,
			); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			if err := rows.Scan(
				&d.EmployeeName,
				&d.TotalOrders,
				&d.AvgDeliveryDays,
				&d.MaxDeliveryDays,
				&d.MinDeliveryDays,
				&d.LateShipments,
				&d.OnTimeRate,
			); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		results = append(results, d)
	}

	c.JSON(http.StatusOK, gin.H{
		"group": group,
		"year":  year,
		"count": len(results),
		"data":  results,
	})
}
