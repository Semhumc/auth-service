package middleware

import (
	"auth-service/internal/models"

	"github.com/gofiber/fiber/v2"
)

func LoginMiddleware(c *fiber.Ctx) error {
	var login models.LoginParams

	if err := c.BodyParser(&login); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if login.Email == "" || login.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "email and password are required",
		})
	}

	c.Locals("login", login)
	return c.Next()
}

func RegisterMiddleware(c *fiber.Ctx) error {
	var register models.RegisterParams

	if err := c.BodyParser(&register); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Zorunlu alanları kontrol et
	if register.Login.Email == "" || register.Login.Password == "" || register.Username == "" || register.Firstname == "" || register.Lastname == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "all fields are required",
		})
	}

	c.Locals("register", register)
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