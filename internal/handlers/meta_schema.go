package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicholasraynes/northwind-mcp/internal/models"
)

// GET /meta/schema
func GetMetaSchema(c *gin.Context) {
	schemas := []models.EndpointSchema{
		{
			Name:        "Health Check",
			Description: "Verifies that the MCP server and database connection are active.",
			Path:        "/health",
			Method:      "GET",
			Parameters:  map[string]string{},
			Returns:     map[string]string{"status": "string"},
		},
		{
			Name:        "Customers",
			Description: "Returns a list of all customers with IDs, company info, and contact details.",
			Path:        "/customers",
			Method:      "GET",
			Parameters:  map[string]string{},
			Returns:     map[string]string{"customer_id": "string", "company_name": "string", "country": "string"},
		},
		{
			Name:        "Orders",
			Description: "Returns order-level details including customer, employee, and dates.",
			Path:        "/orders",
			Method:      "GET",
			Parameters:  map[string]string{"year": "int (optional)"},
			Returns:     map[string]string{"order_id": "int", "order_date": "date", "customer_id": "string"},
		},
		{
			Name:        "Products",
			Description: "Lists all products with category, supplier, and pricing information.",
			Path:        "/products",
			Method:      "GET",
			Parameters:  map[string]string{"category": "string (optional)", "supplier": "string (optional)"},
			Returns:     map[string]string{"product_id": "int", "product_name": "string", "unit_price": "float"},
		},
		{
			Name:        "Suppliers",
			Description: "Lists all suppliers with contact and country info.",
			Path:        "/suppliers",
			Method:      "GET",
			Parameters:  map[string]string{},
			Returns:     map[string]string{"supplier_id": "int", "company_name": "string", "country": "string"},
		},
		{
			Name:        "Order Details",
			Description: "Returns product-level details for each order (quantity, discount, etc).",
			Path:        "/orders/details",
			Method:      "GET",
			Parameters:  map[string]string{"order_id": "int (optional)"},
			Returns:     map[string]string{"product_id": "int", "quantity": "int", "discount": "float"},
		},
		{
			Name:        "Sales by Country",
			Description: "Aggregated total sales per country.",
			Path:        "/summary/sales-by-country",
			Method:      "GET",
			Parameters:  map[string]string{"year": "int (optional)"},
			Returns:     map[string]string{"country": "string", "total_sales": "float"},
		},
		{
			Name:        "Sales by Category",
			Description: "Aggregated total sales grouped by product category.",
			Path:        "/summary/sales-by-category",
			Method:      "GET",
			Parameters:  map[string]string{"year": "int (optional)"},
			Returns:     map[string]string{"category_name": "string", "total_sales": "float"},
		},
		{
			Name:        "Sales by Employee",
			Description: "Aggregated total sales per employee.",
			Path:        "/summary/sales-by-employee",
			Method:      "GET",
			Parameters:  map[string]string{"year": "int (optional)"},
			Returns:     map[string]string{"employee_name": "string", "total_sales": "float"},
		},
		{
			Name:        "Sales by Year",
			Description: "Shows yearly sales totals across all customers.",
			Path:        "/summary/sales-by-year",
			Method:      "GET",
			Parameters:  map[string]string{},
			Returns:     map[string]string{"year": "int", "total_sales": "float"},
		},
		{
			Name:        "Sales by Shipper",
			Description: "Total sales grouped by shipping company.",
			Path:        "/summary/sales-by-shipper",
			Method:      "GET",
			Parameters:  map[string]string{"year": "int (optional)"},
			Returns:     map[string]string{"shipper_name": "string", "total_sales": "float"},
		},
		{
			Name:        "Top Customers",
			Description: "Top customers ranked by total sales revenue.",
			Path:        "/analytics/top-customers",
			Method:      "GET",
			Parameters:  map[string]string{"limit": "int", "country": "string", "year": "int"},
			Returns:     map[string]string{"customer_id": "string", "total_sales": "float"},
		},
		{
			Name:        "Customer Orders",
			Description: "Detailed order history per customer, optionally filtered by year.",
			Path:        "/analytics/customer-orders",
			Method:      "GET",
			Parameters:  map[string]string{"customer_id": "string", "year": "int"},
			Returns:     map[string]string{"order_id": "int", "order_date": "date", "total_amount": "float"},
		},
		{
			Name:        "Customer Lifetime Value",
			Description: "Lifetime revenue and average order per customer.",
			Path:        "/analytics/customer-ltv",
			Method:      "GET",
			Parameters:  map[string]string{"limit": "int", "country": "string"},
			Returns:     map[string]string{"customer_id": "string", "total_sales": "float", "avg_order": "float"},
		},
		{
			Name:        "Customer Retention",
			Description: "Repeat vs new customers and overall retention rate.",
			Path:        "/analytics/customer-retention",
			Method:      "GET",
			Parameters:  map[string]string{"year": "int"},
			Returns:     map[string]string{"retention_rate": "float", "repeat_customers": "int"},
		},
		{
			Name:        "Top Products",
			Description: "Top-selling products with supplier and category info.",
			Path:        "/analytics/top-products",
			Method:      "GET",
			Parameters:  map[string]string{"limit": "int", "year": "int", "category": "string"},
			Returns:     map[string]string{"product_id": "int", "product_name": "string", "total_revenue": "float"},
		},
		{
			Name:        "Supplier Performance",
			Description: "Summarizes supplier sales, category mix, and total product count.",
			Path:        "/analytics/supplier-performance",
			Method:      "GET",
			Parameters:  map[string]string{"limit": "int", "year": "int"},
			Returns:     map[string]string{"supplier_id": "int", "total_revenue": "float", "top_category": "string"},
		},
		{
			Name:        "Inventory Status",
			Description: "Shows stock levels, reorder thresholds, and discontinued flags.",
			Path:        "/analytics/inventory-status",
			Method:      "GET",
			Parameters:  map[string]string{"below_reorder": "bool", "supplier": "string", "category": "string", "discontinued": "bool"},
			Returns:     map[string]string{"product_id": "int", "units_in_stock": "int", "needs_reorder": "bool"},
		},
		{
			Name:        "Employee Performance",
			Description: "Ranks employees by sales and efficiency.",
			Path:        "/analytics/employee-performance",
			Method:      "GET",
			Parameters:  map[string]string{"limit": "int", "year": "int"},
			Returns:     map[string]string{"employee_id": "int", "total_revenue": "float"},
		},
		{
			Name:        "Shipping Costs",
			Description: "Freight cost and order volume per shipping company.",
			Path:        "/analytics/shipping-costs",
			Method:      "GET",
			Parameters:  map[string]string{"limit": "int", "year": "int"},
			Returns:     map[string]string{"shipper_id": "int", "total_freight": "float", "avg_freight": "float"},
		},
		{
			Name:        "Delivery Times",
			Description: "Average, min, and max delivery durations with late shipment counts.",
			Path:        "/analytics/delivery-times",
			Method:      "GET",
			Parameters:  map[string]string{"group": "string ('shipper'|'employee')", "year": "int"},
			Returns:     map[string]string{"avg_delivery_days": "float", "on_time_rate": "float"},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"service":   "Northwind MCP API",
		"version":   "1.0.0",
		"endpoints": schemas,
	})
}
