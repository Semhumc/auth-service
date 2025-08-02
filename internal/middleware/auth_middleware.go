package middleware

import (
	"auth-service/internal/models"
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

// JWT TOKEN VALIDATION MIDDLEWARE - User Management endpoint'leri için
func AuthTokenMiddleware(c *fiber.Ctx) error {
	fmt.Printf("🔐 AuthTokenMiddleware called for path: %s\n", c.Path())
	
	// Token'ı cookie'den veya header'dan al
	token := c.Cookies("access_token")
	if token == "" {
		// Authorization header'dan kontrol et
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			fmt.Printf("❌ No access token provided\n")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "access token required",
			})
		}
		
		// Bearer token formatını kontrol et
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			fmt.Printf("❌ Invalid authorization header format\n")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization header format",
			})
		}
		token = parts[1]
	}

	// Token'ı locals'a kaydet ki handler'da kullanabilelim
	c.Locals("access_token", token)
	fmt.Printf("✅ Token validated and stored in locals\n")
	return c.Next()
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