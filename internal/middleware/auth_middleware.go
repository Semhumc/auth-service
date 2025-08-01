package middleware

import (
	"auth-service/internal/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func LoginMiddleware(c *fiber.Ctx) error {
	fmt.Printf("🔍 LoginMiddleware called\n")
	fmt.Printf("   Method: %s\n", c.Method())
	fmt.Printf("   Path: %s\n", c.Path())
	fmt.Printf("   Content-Type: %s\n", c.Get("Content-Type"))
	
	// Raw body'yi görmek için
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
	fmt.Printf("   Email: %s\n", login.Email)
	fmt.Printf("   Password: %s\n", "***") // Don't log actual password

	if login.Email == "" || login.Password == "" {
		fmt.Printf("❌ Missing required fields\n")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "email and password are required",
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
	
	// Raw body'yi görmek için
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
	fmt.Printf("   Email: %s\n", register.Login.Email)
	fmt.Printf("   Username: %s\n", register.Username)
	fmt.Printf("   Firstname: %s\n", register.Firstname)
	fmt.Printf("   Lastname: %s\n", register.Lastname)
	fmt.Printf("   Password: %s\n", "***") // Don't log actual password

	// Zorunlu alanları kontrol et
	if register.Login.Email == "" || register.Login.Password == "" || register.Username == "" || register.Firstname == "" || register.Lastname == "" {
		fmt.Printf("❌ Missing required fields\n")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "all fields are required",
			"received": fiber.Map{
				"email": register.Login.Email,
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