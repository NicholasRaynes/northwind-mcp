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
		WITH shipper_orders AS (
			SELECT
				s.shipper_id,
				s.company_name,
				o.order_id,
				o.freight,
				c.country
			FROM orders o
			JOIN shippers s ON o.ship_via = s.shipper_id
			JOIN customers c ON c.customer_id = o.customer_id
			WHERE 1=1
`

	args := []any{}
	if year != "" {
		args = append(args, year)
		query += " AND EXTRACT(YEAR FROM o.order_date) = $" + strconv.Itoa(len(args))
	}

	query += `
		),
		shipper_stats AS (
			SELECT
				shipper_id,
				company_name,
				COUNT(DISTINCT order_id) AS total_orders,
				SUM(freight) AS total_freight,
				AVG(freight) AS avg_freight
			FROM shipper_orders
			GROUP BY shipper_id, company_name
		),
		top_destinations AS (
			SELECT
				shipper_id,
				country AS top_destination
			FROM (
				SELECT 
					shipper_id,
					country,
					COUNT(*) AS cnt,
					ROW_NUMBER() OVER (PARTITION BY shipper_id ORDER BY COUNT(*) DESC) AS rn
				FROM shipper_orders
				GROUP BY shipper_id, country
			) x
			WHERE rn = 1
		)
		SELECT 
			ss.shipper_id,
			ss.company_name,
			ss.total_orders,
			ss.total_freight,
			ss.avg_freight,
			td.top_destination
		FROM shipper_stats ss
		LEFT JOIN top_destinations td ON ss.shipper_id = td.shipper_id
		ORDER BY ss.total_freight DESC
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
