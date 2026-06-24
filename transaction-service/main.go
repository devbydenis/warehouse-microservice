package main

import "micro-warehouse/transaction-service/cmd"

// @title Transaction Service API
// @version 1.0
// @description This is a transaction service API for warehouse management system
// @host transaction-service:8085
// @SecurityDefinitions.apiKey Bearer
// @in header
// @name Authorization
// @BasePath /api/v1
func main() {
	cmd.Execute()
}
