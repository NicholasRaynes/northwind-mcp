package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-api/internal/db"
	"github.com/nicholasraynes/northwind-api/internal/models"
)

// GET /analytics/shipping-costs?year=1997&limit=10
func GetShippingCosts(c *gin.Context) {
	limit := 10
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = val
		}
	}
	year := c.Query("year")

	query := `
		WITH shipper_stats AS (
			SELECT
				s.shipper_id,
				s.company_name,
				COUNT(DISTINCT o.order_id) AS total_orders,
				SUM(o.freight) AS total_freight,
				AVG(o.freight) AS avg_freight,
				(
					SELECT c.country
					FROM customers c
					WHERE c.customer_id = o.customer_id
					GROUP BY c.country
					ORDER BY COUNT(*) DESC
					LIMIT 1
				) AS top_destination
			FROM orders o
			JOIN shippers s ON o.ship_via = s.shipper_id
	`

	args := []interface{}{}
	if year != "" {
		query += " WHERE EXTRACT(YEAR FROM o.order_date) = $" + strconv.Itoa(len(args)+1)
		args = append(args, year)
	}

	query += `
			GROUP BY s.shipper_id, s.company_name
		)
		SELECT shipper_id, company_name, total_orders, total_freight, avg_freight, top_destination
		FROM shipper_stats
		ORDER BY total_freight DESC
		LIMIT ` + strconv.Itoa(limit)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	results := []models.ShippingCosts{}
	for rows.Next() {
		var sc models.ShippingCosts
		if err := rows.Scan(
			&sc.ShipperID,
			&sc.CompanyName,
			&sc.TotalOrders,
			&sc.TotalFreight,
			&sc.AvgFreight,
			&sc.TopDestination,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		results = append(results, sc)
	}

	c.JSON(http.StatusOK, gin.H{
		"year":  year,
		"count": len(results),
		"data":  results,
	})
}
