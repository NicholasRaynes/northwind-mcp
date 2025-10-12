package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-mcp/internal/models"
)

// GET /meta/schema
func GetMetaSchema(c *gin.Context) {
	schema := models.APISchema{
		OpenAPI: "3.1.0",
		Info: models.APIInfo{
			Title:       "Northwind MCP API",
			Version:     "1.0.0",
			Description: "Comprehensive self-describing API exposing Northwind analytics, summaries, and data for Copilot Studio integration.",
		},
		Paths: map[string]models.APIPath{
			"/health": {
				Get: &models.APIEndpoint{
					Summary:     "Health Check",
					Description: "Verify that the MCP server and database connection are active.",
					Responses: map[string]models.APIResponse{
						"200": {Description: "Returns a JSON status confirmation."},
					},
				},
			},
			"/customers": {
				Get: &models.APIEndpoint{
					Summary:     "Get Customers",
					Description: "Retrieve all customers with company and contact details.",
					Responses: map[string]models.APIResponse{
						"200": {Description: "List of customers."},
					},
				},
			},
			"/orders": {
				Get: &models.APIEndpoint{
					Summary:     "Get Orders",
					Description: "Retrieve all orders with customer, employee, and date details.",
					Parameters: []models.APIParameter{
						{Name: "year", In: "query", Description: "Filter by year", Required: false, SchemaType: "integer"},
					},
					Responses: map[string]models.APIResponse{
						"200": {Description: "List of orders."},
					},
				},
			},
			"/products": {
				Get: &models.APIEndpoint{
					Summary:     "Get Products",
					Description: "Retrieve product information with supplier and category details.",
					Parameters: []models.APIParameter{
						{Name: "category", In: "query", Description: "Filter by category", Required: false, SchemaType: "string"},
						{Name: "supplier", In: "query", Description: "Filter by supplier", Required: false, SchemaType: "string"},
					},
					Responses: map[string]models.APIResponse{
						"200": {Description: "List of products."},
					},
				},
			},
			"/suppliers": {
				Get: &models.APIEndpoint{
					Summary:     "Get Suppliers",
					Description: "Retrieve all suppliers with contact information.",
					Responses: map[string]models.APIResponse{
						"200": {Description: "List of suppliers."},
					},
				},
			},
			"/orders/details": {
				Get: &models.APIEndpoint{
					Summary:     "Get Order Details",
					Description: "Retrieve product-level details for specific orders.",
					Parameters: []models.APIParameter{
						{Name: "order_id", In: "query", Description: "Filter by order ID", Required: false, SchemaType: "integer"},
					},
					Responses: map[string]models.APIResponse{
						"200": {Description: "Order detail records."},
					},
				},
			},
			"/summary/sales-by-country": {
				Get: &models.APIEndpoint{
					Summary:     "Sales by Country",
					Description: "Total sales aggregated by country.",
					Parameters: []models.APIParameter{
						{Name: "year", In: "query", Description: "Filter by year", Required: false, SchemaType: "integer"},
					},
					Responses: map[string]models.APIResponse{
						"200": {Description: "Sales totals by country."},
					},
				},
			},
			"/summary/sales-by-category": {
				Get: &models.APIEndpoint{
					Summary:     "Sales by Category",
					Description: "Total sales grouped by product category.",
					Parameters: []models.APIParameter{
						{Name: "year", In: "query", Description: "Filter by year", Required: false, SchemaType: "integer"},
					},
					Responses: map[string]models.APIResponse{
						"200": {Description: "Sales totals by category."},
					},
				},
			},
			"/summary/sales-by-employee": {
				Get: &models.APIEndpoint{
					Summary:     "Sales by Employee",
					Description: "Total sales aggregated per employee.",
					Parameters: []models.APIParameter{
						{Name: "year", In: "query", Description: "Filter by year", Required: false, SchemaType: "integer"},
					},
					Responses: map[string]models.APIResponse{
						"200": {Description: "Sales totals by employee."},
					},
				},
			},
			"/summary/sales-by-year": {
				Get: &models.APIEndpoint{
					Summary:     "Sales by Year",
					Description: "Total yearly sales for all orders.",
					Responses: map[string]models.APIResponse{
						"200": {Description: "Annual sales totals."},
					},
				},
			},
			"/summary/sales-by-shipper": {
				Get: &models.APIEndpoint{
					Summary:     "Sales by Shipper",
					Description: "Total sales grouped by shipping company.",
					Parameters: []models.APIParameter{
						{Name: "year", In: "query", Description: "Filter by year", Required: false, SchemaType: "integer"},
					},
					Responses: map[string]models.APIResponse{
						"200": {Description: "Sales totals by shipper."},
					},
				},
			},
			"/analytics/top-customers": {
				Get: &models.APIEndpoint{
					Summary:     "Top Customers",
					Description: "Top customers ranked by total sales revenue.",
					Parameters: []models.APIParameter{
						{Name: "limit", In: "query", Description: "Limit number of results", Required: false, SchemaType: "integer"},
						{Name: "country", In: "query", Description: "Filter by country", Required: false, SchemaType: "string"},
						{Name: "year", In: "query", Description: "Filter by order year", Required: false, SchemaType: "integer"},
					},
					Responses: map[string]models.APIResponse{
						"200": {Description: "Top customers with revenue totals."},
					},
				},
			},
			"/analytics/customer-orders": {
				Get: &models.APIEndpoint{
					Summary:     "Customer Orders",
					Description: "Detailed order history for specific customers.",
					Parameters: []models.APIParameter{
						{Name: "customer_id", In: "query", Description: "Filter by customer ID", Required: false, SchemaType: "string"},
						{Name: "year", In: "query", Description: "Filter by order year", Required: false, SchemaType: "integer"},
					},
					Responses: map[string]models.APIResponse{
						"200": {Description: "List of orders and totals for customer."},
					},
				},
			},
			"/analytics/customer-ltv": {
				Get: &models.APIEndpoint{
					Summary:     "Customer Lifetime Value",
					Description: "Total lifetime revenue and order count per customer.",
					Parameters: []models.APIParameter{
						{Name: "limit", In: "query", Description: "Limit number of results", Required: false, SchemaType: "integer"},
						{Name: "country", In: "query", Description: "Filter by country", Required: false, SchemaType: "string"},
					},
					Responses: map[string]models.APIResponse{
						"200": {Description: "LTV metrics per customer."},
					},
				},
			},
			"/analytics/customer-retention": {
				Get: &models.APIEndpoint{
					Summary:     "Customer Retention",
					Description: "Measures repeat customers and retention rate.",
					Parameters: []models.APIParameter{
						{Name: "year", In: "query", Description: "Filter by active year", Required: false, SchemaType: "integer"},
					},
					Responses: map[string]models.APIResponse{
						"200": {Description: "Retention rate and repeat customer stats."},
					},
				},
			},
			"/analytics/top-products": {
				Get: &models.APIEndpoint{
					Summary:     "Top Products",
					Description: "Top-selling products by total revenue and units sold.",
					Parameters: []models.APIParameter{
						{Name: "limit", In: "query", Description: "Limit number of results", Required: false, SchemaType: "integer"},
						{Name: "year", In: "query", Description: "Filter by year", Required: false, SchemaType: "integer"},
						{Name: "category", In: "query", Description: "Filter by category", Required: false, SchemaType: "string"},
					},
					Responses: map[string]models.APIResponse{
						"200": {Description: "Top product performance results."},
					},
				},
			},
			"/analytics/supplier-performance": {
				Get: &models.APIEndpoint{
					Summary:     "Supplier Performance",
					Description: "Supplier contribution by revenue and top category.",
					Parameters: []models.APIParameter{
						{Name: "limit", In: "query", Description: "Limit results", Required: false, SchemaType: "integer"},
						{Name: "year", In: "query", Description: "Filter by year", Required: false, SchemaType: "integer"},
					},
					Responses: map[string]models.APIResponse{
						"200": {Description: "Supplier performance summary."},
					},
				},
			},
			"/analytics/inventory-status": {
				Get: &models.APIEndpoint{
					Summary:     "Inventory Status",
					Description: "Product stock levels, reorder needs, and discontinued items.",
					Parameters: []models.APIParameter{
						{Name: "below_reorder", In: "query", Description: "Show only low-stock items", Required: false, SchemaType: "boolean"},
						{Name: "supplier", In: "query", Description: "Filter by supplier", Required: false, SchemaType: "string"},
						{Name: "category", In: "query", Description: "Filter by category", Required: false, SchemaType: "string"},
						{Name: "discontinued", In: "query", Description: "Show only discontinued items", Required: false, SchemaType: "boolean"},
					},
					Responses: map[string]models.APIResponse{
						"200": {Description: "Inventory list with stock details."},
					},
				},
			},
			"/analytics/employee-performance": {
				Get: &models.APIEndpoint{
					Summary:     "Employee Performance",
					Description: "Rank employees by sales totals and order count.",
					Parameters: []models.APIParameter{
						{Name: "limit", In: "query", Description: "Limit number of employees", Required: false, SchemaType: "integer"},
						{Name: "year", In: "query", Description: "Filter by order year", Required: false, SchemaType: "integer"},
					},
					Responses: map[string]models.APIResponse{
						"200": {Description: "Employee performance data."},
					},
				},
			},
			"/analytics/shipping-costs": {
				Get: &models.APIEndpoint{
					Summary:     "Shipping Costs",
					Description: "Total freight cost and order volume per shipping company.",
					Parameters: []models.APIParameter{
						{Name: "limit", In: "query", Description: "Limit results", Required: false, SchemaType: "integer"},
						{Name: "year", In: "query", Description: "Filter by year", Required: false, SchemaType: "integer"},
					},
					Responses: map[string]models.APIResponse{
						"200": {Description: "Shipping cost summary."},
					},
				},
			},
			"/analytics/delivery-times": {
				Get: &models.APIEndpoint{
					Summary:     "Delivery Times",
					Description: "Average delivery days, late shipments, and on-time rate by shipper or employee.",
					Parameters: []models.APIParameter{
						{Name: "group", In: "query", Description: "Group by 'shipper' or 'employee'", Required: false, SchemaType: "string"},
						{Name: "year", In: "query", Description: "Filter by year", Required: false, SchemaType: "integer"},
					},
					Responses: map[string]models.APIResponse{
						"200": {Description: "Delivery time and performance metrics."},
					},
				},
			},
		},
	}

	c.JSON(http.StatusOK, schema)
}
