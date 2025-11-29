package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-api/internal/db"
	"github.com/nicholasraynes/northwind-api/internal/models"
)

// GET /analytics/shipping-costs
// Optional parameters: year, shipper_id, company_name
func GetShippingCosts(c *gin.Context) {
	year := c.Query("year")
	shipperID := c.Query("shipper_id")
	companyName := c.Query("company_name")

	args := []any{}
	cteConditions := []string{}

	if year != "" {
		cteConditions = append(cteConditions, fmt.Sprintf("EXTRACT(YEAR FROM o.order_date)::TEXT = $%d", len(args)+1))
		args = append(args, year)
	}

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
	`

	if len(cteConditions) > 0 {
		query += " WHERE " + strings.Join(cteConditions, " AND ")
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
	`

	finalConditions := []string{}

	if shipperID != "" {
		finalConditions = append(finalConditions, fmt.Sprintf("CAST(ss.shipper_id AS TEXT) LIKE $%d", len(args)+1))
		args = append(args, "%"+shipperID+"%")
	}
	if companyName != "" {
		finalConditions = append(finalConditions, fmt.Sprintf("LOWER(ss.company_name) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+companyName+"%")
	}

	if len(finalConditions) > 0 {
		query += " WHERE " + strings.Join(finalConditions, " AND ")
	}

	query += " ORDER BY ss.total_freight DESC"

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

	filters := gin.H{}
	if year != "" {
		filters["year"] = year
	}
	if shipperID != "" {
		filters["shipper_id"] = shipperID
	}
	if companyName != "" {
		filters["company_name"] = companyName
	}

	c.JSON(http.StatusOK, gin.H{
		"filters": filters,
		"count":   len(results),
		"data":    results,
	})
}
