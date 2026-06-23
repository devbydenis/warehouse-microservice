package main

import "micro-warehouse/product-service/cmd"

// @title Product Service API
// @version 2.0
// @description This is the API documentation for the Product Service.
// @host product-service:8082
// @BasePath /api/v1
func main() {
	cmd.Execute()
}