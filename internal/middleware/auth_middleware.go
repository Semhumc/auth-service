package middleware

import (
	"auth-service/internal/models"
	"auth-service/internal/services"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func LoginMiddleware(c *fiber.Ctx) error {
	fmt.Printf("🔍 LoginMiddleware called\n")
	fmt.Printf("   Method: %s\n", c.Method())
	fmt.Printf("   Path: %s\n", c.Path())
	fmt.Printf("   Content-Type: %s\n", c.Get("Content-Type"))
	
	rawBody := c.Body()
	fmt.Printf("   Raw Body: %s\n", string(rawBody))

	var login models.LoginParams

	if err := c.BodyParser(&login); err != nil {
		fmt.Printf("❌ Body parsing failed: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
			"details": err.Error(),
		})
	}

	fmt.Printf("✅ Parsed login data:\n")
	fmt.Printf("   Username: %s\n", login.Username)
	fmt.Printf("   Password: %s\n", "***")

	if login.Username == "" || login.Password == "" {
		fmt.Printf("❌ Missing required fields\n")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "username and password are required",
		})
	}

	c.Locals("login", login)
	fmt.Printf("✅ LoginMiddleware completed successfully\n")
	return c.Next()
}

func RegisterMiddleware(c *fiber.Ctx) error {
	fmt.Printf("🔍 RegisterMiddleware called\n")
	fmt.Printf("   Method: %s\n", c.Method())
	fmt.Printf("   Path: %s\n", c.Path())
	fmt.Printf("   Content-Type: %s\n", c.Get("Content-Type"))
	
	rawBody := c.Body()
	fmt.Printf("   Raw Body: %s\n", string(rawBody))

	var register models.RegisterParams

	if err := c.BodyParser(&register); err != nil {
		fmt.Printf("❌ Body parsing failed: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
			"details": err.Error(),
		})
	}

	fmt.Printf("✅ Parsed register data:\n")
	fmt.Printf("   Email: %s\n", register.Email)
	fmt.Printf("   Username: %s\n", register.Username)
	fmt.Printf("   Firstname: %s\n", register.Firstname)
	fmt.Printf("   Lastname: %s\n", register.Lastname)
	fmt.Printf("   Password: %s\n", "***")

	if register.Email == "" || register.Password == "" || register.Username == "" || register.Firstname == "" || register.Lastname == "" {
		fmt.Printf("❌ Missing required fields\n")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "all fields are required",
			"received": fiber.Map{
				"email": register.Email,
				"username": register.Username,
				"firstname": register.Firstname,
				"lastname": register.Lastname,
				"password": "***",
			},
		})
	}

	c.Locals("register", register)
	fmt.Printf("✅ RegisterMiddleware completed successfully\n")
	return c.Next()
}

func NewAuthTokenMiddleware(keycloakService *services.KeycloakService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Printf("🔐 AuthTokenMiddleware called for path: %s\n", c.Path())

		// 1. Get access token from header or cookie
		accessToken := c.Cookies("access_token")
		if accessToken == "" {
			authHeader := c.Get("Authorization")
			if authHeader == "" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "access token required"})
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid authorization header format"})
			}
			accessToken = parts[1]
		}

		// 2. Introspect the token
		ctx := c.Context()
		result, err := keycloakService.Gocloak.RetrospectToken(ctx, accessToken, keycloakService.ClientId, keycloakService.ClientSecret, keycloakService.Realm)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to introspect token", "details": err.Error()})
		}

		// 3. Check if token is active
		if !*result.Active {
			fmt.Println("Token is inactive, attempting refresh")

			// 3a. Get refresh token from cookie
			refreshToken := c.Cookies("refresh_token")
			if refreshToken == "" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "session expired, no refresh token found"})
			}

			// 3b. Attempt to refresh the token
			newTokens, err := keycloakService.RefreshToken(refreshToken)
			if err != nil {
				// Clear cookies if refresh fails
				c.Cookie(&fiber.Cookie{Name: "access_token", MaxAge: -1})
				c.Cookie(&fiber.Cookie{Name: "refresh_token", MaxAge: -1})
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "session expired, token refresh failed", "details": err.Error()})
			}

			// 3c. Refresh successful, set new tokens in cookies
			c.Cookie(&fiber.Cookie{
				Name:     "access_token",
				Value:    newTokens.AccessToken,
				HTTPOnly: true,
				Secure:   false,
				SameSite: "Lax",
			})
			c.Cookie(&fiber.Cookie{
				Name:     "refresh_token",
				Value:    newTokens.RefreshToken,
				HTTPOnly: true,
				Secure:   false,
				SameSite: "Lax",
			})

			// Use the new access token for the current request
			c.Locals("access_token", newTokens.AccessToken)
			fmt.Println("✅ Token refreshed successfully")
			return c.Next()
		}

		// 4. Token is active, proceed
		c.Locals("access_token", accessToken)
		fmt.Printf("✅ Token validated and stored in locals\n")
		return c.Next()
	}
}

func GetUserMiddleware(c *fiber.Ctx) error {
	userID := c.Params("id")

	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user ID is required",
		})
	}

	c.Locals("userID", userID)
	return c.Next()
}

func UpdateMiddleware(c *fiber.Ctx) error {
	userID := c.Params("id")

	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user ID is required",
		})
	}

	c.Locals("userID", userID)
	return c.Next()
}

func DeleteMiddleware(c *fiber.Ctx) error {
	userID := c.Params("id")

	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user ID is required",
		})
	}

	c.Locals("userID", userID)
	return c.Next()
}
