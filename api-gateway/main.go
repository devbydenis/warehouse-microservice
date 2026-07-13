package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	jwtConf "micro-warehouse/api-gateway/configs"
	"micro-warehouse/api-gateway/controller"
	swaggerAgg "micro-warehouse/api-gateway/internal/swagger"
	"micro-warehouse/api-gateway/middleware"

	_ "micro-warehouse/api-gateway/docs"
)

type ServiceConfig struct {
	Name string
	URL  string
}

type Config struct {
	Services map[string]ServiceConfig
	Port     string
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: .env file not found, trying env.example")
		if err := godotenv.Load("env.example"); err != nil {
			log.Println("Warning: env.example file not found, using default configuration")
		}
	}

	config := loadConfig()
	jwtConfig := jwtConf.LoadJWTConfig()
	redisConfig := jwtConf.LoadRedisConfig()

	app := fiber.New(fiber.Config{
		AppName:      "Warehouse Project API Gateway",
		ServerHeader: "Warehouse-API-Gateway",
	})

	rateLimiter := middleware.DefaultRateLimiterConfig()

	var redisClient *redis.Client
	if redisConfig.Host != "" {
		redisClient = jwtConf.NewRedisClient(redisConfig)
	}
	redisRateLimiterConfig := middleware.RedisRateLimiterConfig{
		Max:         rateLimiter.Max,
		Expiration:  rateLimiter.Expiration,
		RedisClient: redisClient,
	}

	app.Use(middleware.RedisGlobalRateLimiter(redisRateLimiterConfig))
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "OK",
			"message": "API Gateway is running",
			"services": fiber.Map{
				"user-service":         config.Services["user"].URL,
				"product-service":      config.Services["product"].URL,
				"merchant-service":     config.Services["merchant"].URL,
				"notification-service": config.Services["notification"].URL,
				"transaction-service":  config.Services["transaction"].URL,
				"midtrans-service":     config.Services["midtrans"].URL,
				"warehouse-service":    config.Services["warehouse"].URL,
				"role-service":         config.Services["role"].URL,
				"assign-role-service":  config.Services["assign-role"].URL,
				// "auth-service":         config.Services["auth"].URL,
			},
		})
	})

	authController := controller.NewAuthController(config.Services["auth"].URL, jwtConfig)

	setupAuthRoutes(app, authController, redisRateLimiterConfig)
	setupMidtransCallbackRoutes(app, config.Services["midtrans"])
	setupProtectedRoutes(app, config, jwtConfig, redisRateLimiterConfig)
	setupSwaggerRoutes(app)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Not Found",
			"message": "Route not found",
			"path":    c.Path(),
		})
	})

	log.Printf("API Gateway starting on port: %s", config.Port)
	log.Fatal(app.Listen(":" + config.Port))
}

func loadConfig() Config {
	config := Config{
		Port: getEnv("PORT", "8080"),
		Services: map[string]ServiceConfig{
			"user": {
				Name: "user-service",
				URL:  getEnv("USER_SERVICE_URL", "http://localhost:8081"),
			},
			"role": {
				Name: "role-service",
				URL:  getEnv("USER_SERVICE_URL", "http://localhost:8081"),
			},
			"assign-role": {
				Name: "assign-role-service",
				URL:  getEnv("USER_SERVICE_URL", "http://localhost:8081"),
			},
			"auth": {
				Name: "auth-service",
				URL:  getEnv("USER_SERVICE_URL", "http://localhost:8081"),
			},
			"product": {
				Name: "product-service",
				URL:  getEnv("PRODUCT_SERVICE_URL", "http://localhost:8082"),
			},
			"merchant": {
				Name: "merchant-service",
				URL:  getEnv("MERCHANT_SERVICE_URL", "http://localhost:8084"),
			},
			"notification": {
				Name: "notification-service",
				URL:  getEnv("NOTIFICATION_SERVICE_URL", "http://localhost:8086"),
			},
			"transaction": {
				Name: "transaction-service",
				URL:  getEnv("TRANSACTION_SERVICE_URL", "http://localhost:8085"),
			},
			"midtrans": {
				Name: "midtrans-service",
				URL:  getEnv("TRANSACTION_SERVICE_URL", "http://localhost:8085"),
			},
			"warehouse": {
				Name: "warehouse-service",
				URL:  getEnv("WAREHOUSE_SERVICE_URL", "http://localhost:8083"),
			},
		},
	}

	log.Println("Service Configuration")
	for name, service := range config.Services {
		log.Printf("Service: %s, URL: %s", name, service.URL)
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	
	return defaultValue
}

func proxyRequestWithPath(c *fiber.Ctx, targetURL string, basePath string) error {
	fullPath := c.Path()
	fullURL := targetURL + fullPath

	queryParams := c.Context().QueryArgs().String()
	if queryParams != "" {
		fullURL += "?" + queryParams
	}

	client := &http.Client{}
	body := c.Body()

	req, err := http.NewRequest(c.Method(), fullURL, bytes.NewReader(body))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": "Failed to create request",
		})
	}

	for key, values := range c.GetReqHeaders() {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	req.Header.Set("X-Gateway", "warehouse-api-gateway")
	req.Header.Set("X-Internal-Request", "true")

	addUserHeaders(req, c)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request to %s: %v", fullURL, err)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   "Bad Gateway",
			"message": "Service unavailable",
			"service": targetURL,
		})
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": "Failed to read response body",
		})
	}

	for key, values := range resp.Header {
		for _, value := range values {
			c.Set(key, value)
		}
	}

	return c.Status(resp.StatusCode).Send(respBody)
}

func proxyRequest(c *fiber.Ctx, targetURL string) error {
	path := c.Params("*")
	if path == "" {
		path = c.Path()
	}

	path = strings.TrimPrefix(path, "api/v1")

	fullURL := targetURL
	if !strings.HasSuffix(targetURL, "/") && !strings.HasPrefix(path, "/") {
		fullURL += "/"
	}
	fullURL += path

	queryParams := c.Context().QueryArgs().String()
	if queryParams != "" {
		fullURL += "?" + queryParams
	}

	client := &http.Client{}
	body := c.Body()

	req, err := http.NewRequest(c.Method(), fullURL, bytes.NewReader(body))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": "Failed to create request",
		})
	}

	for key, values := range c.GetReqHeaders() {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	req.Header.Set("X-Gateway", "warehouse-api-gateway")
	req.Header.Set("X-Internal-Request", "true")

	addUserHeaders(req, c)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request to %s: %v", fullURL, err)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   "Bad Gateway",
			"message": "Service unavailable",
			"service": targetURL,
		})
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": "Failed to read response body",
		})
	}

	for key, values := range resp.Header {
		for _, value := range values {
			c.Set(key, value)
		}
	}

	return c.Status(resp.StatusCode).Send(respBody)
}

func addUserHeaders(req *http.Request, c *fiber.Ctx) {
	if userID := c.Locals("user_id"); userID != nil {
		req.Header.Set("X-User-ID", fmt.Sprintf("%v", userID))
	}
	if userEmail := c.Locals("user_email"); userEmail != nil {
		req.Header.Set("X-User-Email", fmt.Sprintf("%v", userEmail))
	}
	if userRoles := c.Locals("user_roles"); userRoles != nil {
		req.Header.Set("X-User-Roles", fmt.Sprintf("%v", userRoles))
	}
}

func setupSwaggerRoutes(app *fiber.App) {
	// swagger ui - kita serve dari gateway
	app.Get("/swagger", func(c *fiber.Ctx) error {
		html := `<!DOCTYPE html>
	<html>
		<head><title>API Docs</title>
			<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css">
		</head>
		<body>
			<div id="swagger-ui"></div>
			<script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
			<script>
				SwaggerUIBundle({ url: "/swagger/aggregated.json", dom_id: '#swagger-ui', tagsSorter: 'alpha' })
			</script>
		</body>
	</html>`

		c.Set("Content-Type", "text/html")
		c.Set("X-Gateway", "warehouse-api-gateway")
		c.Set("X-Internal-Request", "true")

		return c.SendString(html)
	})

	// Endpoint aggregated spec
	app.Get("/swagger/aggregated.json", func(c *fiber.Ctx) error {
		gatewayURL := c.BaseURL()
		spec, err := swaggerAgg.AggregateSpecs(gatewayURL)
		if err != nil {
			log.Printf("[SetupSwaggerRoutes] AggregateSpecs: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to aggregate specs",
			})
		}

		c.Status(fiber.StatusOK)
		c.Set("Content-Type", "application/json")

		data, err := json.Marshal(spec)
		if err != nil {
			log.Printf("[SetupSwaggerRoutes] Marshal: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to marshal specs",
			})
		}

		return c.Send(data)
	})
}

func setupAuthRoutes(app *fiber.App, authController *controller.AuthController, rateLimiterConfig middleware.RedisRateLimiterConfig) {
	authGroup := app.Group("/api/v1/auth")

	authGroup.Use(middleware.RedisAuthRateLimiter(rateLimiterConfig))
	authGroup.Post("/login", authController.Login)
}

func setupMidtransCallbackRoutes(app *fiber.App, service ServiceConfig) {
	app.Post("/api/v1/midtrans/callback", func(c *fiber.Ctx) error {
		client := &http.Client{}
		fullURL := service.URL + "/api/v1/midtrans/callback"
		body := c.Body()

		req, err := http.NewRequest(c.Method(), fullURL, bytes.NewReader(body))
		if err != nil {
			log.Printf("Error creating request: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Internal Server Error",
				"message": "Failed to create request",
			})
		}

		for key, values := range c.GetReqHeaders() {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		req.Header.Set("X-Gateway", "warehouse-api-gateway")
		req.Header.Set("X-Internal-Request", "true")

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error making request to %s: %v", fullURL, err)
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"error":   "Bad Gateway",
				"message": "Service unavailable",
				"service": service.URL,
			})
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Internal Server Error",
				"message": "Failed to read response body",
			})
		}

		for key, values := range resp.Header {
			for _, value := range values {
				c.Set(key, value)
			}
		}

		return c.Status(resp.StatusCode).Send(respBody)
	})
}

func setupProtectedRoutes(app *fiber.App, config Config, jwtConfig middleware.JWTConfig, rateLimiterConfig middleware.RedisRateLimiterConfig) {
	protectedGroup := app.Group("/api/v1", middleware.JWTAuthMiddleware(jwtConfig))

	protectedGroup.Use(middleware.RedisAPIRateLimiter(rateLimiterConfig))

	setupUserRoutes(protectedGroup, config.Services["user"])
	setupRoleRoutes(protectedGroup, config.Services["role"])
	setupAssignRoleRoutes(protectedGroup, config.Services["assign-role"])
	setupProductRoutes(protectedGroup, config.Services["product"])
	setupMerchantRoutes(protectedGroup, config.Services["merchant"])
	setupTransactionRoutes(protectedGroup, config.Services["transaction"])
	setupWarehouseRoutes(protectedGroup, config.Services["warehouse"])
}

func setupUserRoutes(router fiber.Router, service ServiceConfig) {
	userGroup := router.Group("/users")

	userGroup.All("/*", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/users")
	})
	userGroup.All("/", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/users")
	})

	uploadUserGroup := router.Group("/upload")

	uploadUserGroup.All("/*", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/upload")
	})
	uploadUserGroup.All("/", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/upload")
	})
}

func setupRoleRoutes(router fiber.Router, service ServiceConfig) {
	roleGroup := router.Group("/roles")

	roleGroup.All("/*", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/roles")
	})

	roleGroup.All("/", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/roles")
	})
}

func setupAssignRoleRoutes(router fiber.Router, service ServiceConfig) {
	assignRoleGroup := router.Group("/assign-role")

	assignRoleGroup.All("/*", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/assign-role")
	})

	assignRoleGroup.All("/", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/assign-role")
	})
}

func setupProductRoutes(router fiber.Router, service ServiceConfig) {
	productGroup := router.Group("/products")

	productGroup.All("/*", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/products")
	})
	productGroup.All("/", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/products")
	})

	categoryGroup := router.Group("/categories")

	categoryGroup.All("/*", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/categories")
	})
	categoryGroup.All("/", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/categories")
	})

	uploadGroup := router.Group("/upload-product")

	uploadGroup.All("/*", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/upload-product")
	})
	uploadGroup.All("/", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/upload-product")
	})
}

func setupMerchantRoutes(router fiber.Router, service ServiceConfig) {
	merchantGroup := router.Group("/merchants")

	merchantGroup.All("/*", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/merchants")
	})
	merchantGroup.All("/", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/merchants")
	})

	merchantProductGroup := router.Group("/merchant-products")

	merchantProductGroup.All("/*", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/merchant-products")
	})
	merchantProductGroup.All("/", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/merchant-products")
	})

	uploadMerchantProductGroup := router.Group("/upload-merchant")

	uploadMerchantProductGroup.All("/*", func(c *fiber.Ctx) error {
		return proxyRequest(c, service.URL)
	})
	uploadMerchantProductGroup.All("/", func(c *fiber.Ctx) error {
		return proxyRequest(c, service.URL)
	})
}

func setupTransactionRoutes(router fiber.Router, service ServiceConfig) {
	transactionGroup := router.Group("/transactions")

	transactionGroup.All("/*", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/transactions")
	})
	transactionGroup.All("/", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/transactions")
	})

	dashboardGroup := router.Group("/dashboard")

	dashboardGroup.All("/*", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/dashboard")
	})
	dashboardGroup.All("/", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/dashboard")
	})
}

func setupWarehouseRoutes(router fiber.Router, service ServiceConfig) {
	warehouseGroup := router.Group("/warehouses")

	warehouseGroup.All("/*", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/warehouses")
	})
	warehouseGroup.All("/", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/warehouses")
	})

	warehouseProductGroup := router.Group("/warehouse-products")

	warehouseProductGroup.All("/*", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/warehouse-products")
	})
	warehouseProductGroup.All("/", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/warehouse-products")
	})

	uploadWarehouseGroup := router.Group("/upload-warehouse")

	uploadWarehouseGroup.All("/*", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/upload-warehouse")
	})
	uploadWarehouseGroup.All("/", func(c *fiber.Ctx) error {
		return proxyRequestWithPath(c, service.URL, "/api/v1/upload-warehouse")
	})
}
