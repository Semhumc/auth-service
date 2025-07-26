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

	if register.Login.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "password required",
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

	var userPayload models.UserPayload

	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user ID is required",
		})
	}

	if err := c.BodyParser(&userPayload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	c.Locals("userID",userID)
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
