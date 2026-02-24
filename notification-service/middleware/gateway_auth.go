package middleware

import "github.com/gofiber/fiber/v2"

func GatewayAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		internalRequest := c.Get("X-Internal-Request")
		gatewayHeader := c.Get("X-Gateway")

		if internalRequest != "true" || gatewayHeader != "warehouse-api-gateway" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   "Forbidden",
				"message": "Direct access to service is not allowed. Please use API Gateway.",
				"code":    "DIRECT_ACCESS_FORBIDDEN",
			})
		}

		return c.Next()
	}
}

func OptionalGatewayAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		internalRequest := c.Get("X-Internal-Request")
		gatewayHeader := c.Get("X-Gateway")

		if internalRequest != "true" || gatewayHeader != "warehouse-api-gateway" {
			return c.Next()
		}

		return c.Next()
	}
}
