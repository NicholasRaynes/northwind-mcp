package main

import (
	"fmt"
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

	fmt.Printf("🚀 Server running on port %s\n", port)
	err = r.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}
