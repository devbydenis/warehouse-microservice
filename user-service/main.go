package main

import "micro-warehouse/user-service/cmd"

// @title User Service API
// @version 1.0
// @description This is the API documentation for the User Service.
// @host user-service:8081
// @BasePath /api/v1
func main() {
	cmd.Execute()
}