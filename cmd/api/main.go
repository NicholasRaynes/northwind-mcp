package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/nicholasraynes/northwind-mcp/internal/db"
	"github.com/nicholasraynes/northwind-mcp/internal/handlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db.Connect()

	r := gin.Default()
	r.GET("/health", handlers.Health)
	r.GET("/customers", handlers.GetCustomers)
	r.GET("/orders", handlers.GetOrders)
	r.GET("/products", handlers.GetProducts)
	r.GET("/suppliers", handlers.GetSuppliers)
	r.GET("/orders/details", handlers.GetOrderDetails)
	r.GET("/summary/sales-by-country", handlers.GetSalesByCountry)
	r.GET("/summary/sales-by-category", handlers.GetSalesByCategory)
	r.GET("/summary/sales-by-employee", handlers.GetSalesByEmployee)
	r.GET("/summary/sales-by-year", handlers.GetSalesByYear)
	r.GET("/summary/sales-by-shipper", handlers.GetSalesByShipper)
	r.GET("/analytics/top-customers", handlers.GetTopCustomers)
	r.GET("/analytics/customer-orders", handlers.GetCustomerOrders)
	r.GET("/analytics/customer-ltv", handlers.GetCustomerLTV)
	r.GET("/analytics/customer-retention", handlers.GetCustomerRetention)
	r.GET("/analytics/top-products", handlers.GetTopProducts)
	r.GET("/analytics/supplier-performance", handlers.GetSupplierPerformance)
	r.GET("/analytics/inventory-status", handlers.GetInventoryStatus)
	r.GET("/analytics/employee-performance", handlers.GetEmployeePerformance)
	r.GET("/analytics/shipping-costs", handlers.GetShippingCosts)
	r.GET("/analytics/delivery-times", handlers.GetDeliveryTimes)
	r.StaticFile("/openapi.json", "./schema/openapi.json")

	err = r.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}
