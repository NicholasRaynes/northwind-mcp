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
			Description: "Comprehensive OpenAPI description for Copilot Studio integration. Provides analytics, summaries, and core Northwind data endpoints with contextual metadata for LLM-based querying.",
		},
		Tags: []models.APITag{
			{Name: "Core", Description: "Base customer, order, and product endpoints for data retrieval."},
			{Name: "Orders", Description: "Endpoints related to order details, customers, employees, and shippers."},
			{Name: "Products", Description: "Endpoints related to products, suppliers, and inventory data."},
			{Name: "Summaries", Description: "Sales aggregations by country, category, employee, and shipper."},
			{Name: "Analytics", Description: "Advanced KPIs, rankings, and performance-based insights for customers, suppliers, and employees."},
			{Name: "Performance", Description: "Operational analytics related to employees, shipping, and delivery performance."},
			{Name: "Financial", Description: "Endpoints focusing on revenue, lifetime value, and cost analysis."},
		},
		Paths: map[string]models.APIPath{
			"/health": {
				Get: &models.APIEndpoint{
					Tags:        []string{"Core"},
					OperationID: "checkHealth",
					Summary:     "Health Check",
					Description: "Verify MCP server and database connection.",
					Responses: map[string]models.APIResponse{"200": {Description: "{"status":"ok"}"}},
				},
			},
			"/customers": {
				Get: &models.APIEndpoint{
					Tags:        []string{"Core", "Financial"},
					OperationID: "getCustomers",
					Summary:     "Get Customers",
					Description: "Retrieve all customers, including company names, locations, and contact details.",
					Responses: map[string]models.APIResponse{"200": {Description: "[{"customer_id":"ALFKI","company_name":"Alfreds Futterkiste","country":"Germany"}]"}},
				},
			},
			"/orders": {
				Get: &models.APIEndpoint{
					Tags:        []string{"Core", "Orders", "Financial"},
					OperationID: "getOrders",
					Summary:     "Get Orders",
					Description: "Retrieve all orders with customer, employee, shipper, and date details.",
					Parameters: []models.APIParameter{{Name: "year", In: "query", SchemaType: "integer"}},
					Responses: map[string]models.APIResponse{"200": {Description: "[{"order_id":10248,"customer_id":"VINET","order_date":"1996-07-04"}]"}},
				},
			},
			"/products": {
				Get: &models.APIEndpoint{
					Tags:        []string{"Core", "Products"},
					OperationID: "getProducts",
					Summary:     "Get Products",
					Description: "Retrieve product information with supplier, price, and stock details.",
					Parameters: []models.APIParameter{{Name: "category", In: "query", SchemaType: "string"}, {Name: "supplier", In: "query", SchemaType: "string"}},
					Responses: map[string]models.APIResponse{"200": {Description: "[{"product_id":1,"product_name":"Chai","category":"Beverages"}]"}},
				},
			},
			"/suppliers": {
				Get: &models.APIEndpoint{
					Tags:        []string{"Products", "Analytics"},
					OperationID: "getSuppliers",
					Summary:     "Get Suppliers",
					Description: "Retrieve all supplier data, including company names, contact info, and countries.",
					Responses: map[string]models.APIResponse{"200": {Description: "[{"supplier_id":1,"company_name":"Exotic Liquids","country":"UK"}]"}},
				},
			},
			"/orders/details": {
				Get: &models.APIEndpoint{
					Tags:        []string{"Orders", "Products"},
					OperationID: "getOrderDetails",
					Summary:     "Get Order Details",
					Description: "Retrieve product-level order details including quantity, price, and discounts.",
					Parameters: []models.APIParameter{{Name: "order_id", In: "query", SchemaType: "integer"}},
					Responses: map[string]models.APIResponse{"200": {Description: "[{"order_id":10248,"product_id":11,"unit_price":14.00,"quantity":12}]"}},
				},
			},
			"/summary/sales-by-country": {
				Get: &models.APIEndpoint{
					Tags:        []string{"Summaries", "Financial"},
					OperationID: "getSalesByCountry",
					Summary:     "Sales by Country",
					Description: "Aggregate total sales grouped by customer country.",
					Responses: map[string]models.APIResponse{"200": {Description: "[{"country":"USA","total_sales":12345.67}]"}},
				},
			},
			"/summary/sales-by-category": {
				Get: &models.APIEndpoint{
					Tags:        []string{"Summaries", "Products"},
					OperationID: "getSalesByCategory",
					Summary:     "Sales by Category",
					Description: "Total revenue aggregated by product category.",
					Responses: map[string]models.APIResponse{"200": {Description: "[{"category":"Beverages","total_sales":23456.78}]"}},
				},
			},
			"/summary/sales-by-employee": {
				Get: &models.APIEndpoint{
					Tags:        []string{"Summaries", "Performance"},
					OperationID: "getSalesByEmployee",
					Summary:     "Sales by Employee",
					Description: "Total sales and order counts grouped by employee.",
					Responses: map[string]models.APIResponse{"200": {Description: "[{"employee":"Nancy Davolio","total_sales":89000.5}]"}},
				},
			},
			"/summary/sales-by-year": {
				Get: &models.APIEndpoint{
					Tags:        []string{"Summaries", "Financial"},
					OperationID: "getSalesByYear",
					Summary:     "Sales by Year",
					Description: "Annual sales totals across all countries and categories.",
					Responses: map[string]models.APIResponse{"200": {Description: "[{"year":1997,"total_sales":345678.9}]"}},
				},
			},
			"/summary/sales-by-shipper": {
				Get: &models.APIEndpoint{
					Tags:        []string{"Summaries", "Orders"},
					OperationID: "getSalesByShipper",
					Summary:     "Sales by Shipper",
					Description: "Total sales volume by freight company.",
					Responses: map[string]models.APIResponse{"200": {Description: "[{"shipper":"Speedy Express","total_sales":56789.12}]"}},
				},
			},
			"/analytics/top-customers": {
				Get: &models.APIEndpoint{
					Tags:        []string{"Analytics", "Financial"},
					OperationID: "getTopCustomers",
					Summary:     "Top Customers",
					Description: "Retrieve top customers ranked by total revenue.",
					Responses: map[string]models.APIResponse{"200": {Description: "[{"customer":"Alfreds Futterkiste","total_sales":50000}]"}},
				},
			},
			"/analytics/customer-orders": {
				Get: &models.APIEndpoint{
					Tags:        []string{"Analytics", "Orders"},
					OperationID: "getCustomerOrders",
					Summary:     "Customer Orders",
					Description: "Detailed order history per customer for timeline or frequency analysis.",
					Responses: map[string]models.APIResponse{"200": {Description: "[{"customer_id":"VINET","orders":15}]"}},
				},
			},
			"/analytics/customer-ltv": {
				Get: &models.APIEndpoint{
					Tags:        []string{"Analytics", "Financial"},
					OperationID: "getCustomerLTV",
					Summary:     "Customer Lifetime Value",
					Description: "Calculate total lifetime value and average revenue per customer.",
					Responses: map[string]models.APIResponse{"200": {Description: "[{"customer":"Around the Horn","ltv":45000}]"}},
				},
			},
			"/analytics/customer-retention": {
				Get: &models.APIEndpoint{
					Tags:        []string{"Analytics", "Performance"},
					OperationID: "getCustomerRetention",
					Summary:     "Customer Retention",
					Description: "Measures repeat customers and year-over-year retention.",
					Responses: map[string]models.APIResponse{"200": {Description: "[{"year":1998,"retention_rate":0.82}]"}},
				},
			},
			"/analytics/top-products": {
				Get: &models.APIEndpoint{
					Tags:        []string{"Analytics", "Products"},
					OperationID: "getTopProducts",
					Summary:     "Top Products",
					Description: "Retrieve top-selling products by total sales and quantity.",
					Responses: map[string]models.APIResponse{"200": {Description: "[{"product":"Chai","sales":20000.5}]"}},
				},
			},
			"/analytics/supplier-performance": {
				Get: &models.APIEndpoint{
					Tags:        []string{"Analytics", "Products", "Performance"},
					OperationID: "getSupplierPerformance",
					Summary:     "Supplier Performance",
					Description: "Ranks suppliers by total revenue and delivery efficiency.",
					Responses: map[string]models.APIResponse{"200": {Description: "[{"supplier":"Exotic Liquids","sales":12345.67}]"}},
				},
			},
			"/analytics/inventory-status": {
				Get: &models.APIEndpoint{
					Tags:        []string{"Analytics", "Products", "Performance"},
					OperationID: "getInventoryStatus",
					Summary:     "Inventory Status",
					Description: "View stock levels, reorder needs, and discontinued items.",
					Responses: map[string]models.APIResponse{"200": {Description: "[{"product":"Chai","units_in_stock":39,"discontinued":false}]"}},
				},
			},
			"/analytics/employee-performance": {
				Get: &models.APIEndpoint{
					Tags:        []string{"Analytics", "Performance"},
					OperationID: "getEmployeePerformance",
					Summary:     "Employee Performance",
					Description: "Ranks employees based on sales totals and fulfilled orders.",
					Responses: map[string]models.APIResponse{"200": {Description: "[{"employee":"Nancy Davolio","total_sales":89000.5}]"}},
				},
			},
			"/analytics/shipping-costs": {
				Get: &models.APIEndpoint{
					Tags:        []string{"Analytics", "Orders", "Financial"},
					OperationID: "getShippingCosts",
					Summary:     "Shipping Costs",
					Description: "Shows freight costs and order volumes per shipper.",
					Responses: map[string]models.APIResponse{"200": {Description: "[{"shipper":"United Package","total_freight":12000}]"}},
				},
			},
			"/analytics/delivery-times": {
				Get: &models.APIEndpoint{
					Tags:        []string{"Analytics", "Performance", "Orders"},
					OperationID: "getDeliveryTimes",
					Summary:     "Delivery Times",
					Description: "Average delivery durations, delays, and on-time shipment rates.",
					Responses: map[string]models.APIResponse{"200": {Description: "[{"shipper":"Speedy Express","avg_days":3.5,"on_time_rate":0.92}]"}},
				},
			},
		},
	}

	c.JSON(http.StatusOK, schema)
}
