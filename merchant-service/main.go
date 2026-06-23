package main

import "micro-warehouse/merchant-service/cmd"

// @title Merchant Service API
// @version 1.0
// @description This is the API for the merchant service
// @host merchant-service:8084
// @BasePath /api/v1
func main() {
	cmd.Execute()
}
