package main

import "micro-warehouse/warehouse-service/cmd"

// @title Warehouse Service API
// @version 1.0
// @description This is the API documentation for the Warehouse Service.
// @host warehouse-service:8083
// @SecurityDefinitions.apiKey Bearer
// @in header
// @name Authorization
// @BasePath /api/v1
func main() {
	cmd.Execute()
}
