package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Roles  string `json:"roles"`
	jwt.RegisteredClaims
}

type JWTConfig struct {
	SecretKey string
	Issuer    string
	Duration  time.Duration
}

func JWTAuthMiddleware(config JWTConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if isPublicRoute(c.Path()) {
			return c.Next()
		}

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "Authorization header required",
			})
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "Invalid token format. Use 'Bearer <token>'",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := validateJWT(tokenString, config.SecretKey)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": err.Error(),
			})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("roles", claims.Roles)

		return c.Next()
	}
}

func RoleAuthMiddlware(requiredRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoles := c.Locals("roles")
		if userRoles == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "User context not found",
			})
		}

		roles, ok := userRoles.([]string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "Invalid user roles format",
			})
		}

		// check if user has required roles
		for _, requiredRole := range requiredRoles {
			for _, userRole := range roles {
				if userRole == requiredRole {
					return c.Next()
				}
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   "Forbidden",
			"message": "User does not have required roles",
		})
	}
}

func validateJWT(tokenString, secretKey string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

func GenerateJWT(userID uint, email string, roles string, config JWTConfig) (string, error) {
	claims := &JWTClaims{
		UserID: userID,
		Email:  email,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.Duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    config.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.SecretKey))
}

func isPublicRoute(path string) bool {
	publicRoutes := []string{
		"/health",
		"/api/v1/auth/login",
		"/api/v1/midtrans/callback",
	}

	for _, route := range publicRoutes {
		if strings.HasPrefix(path, route) {
			return true
		}
	}

	return false
}
