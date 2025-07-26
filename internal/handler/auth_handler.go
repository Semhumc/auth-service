package handler

import (
	"auth-service/internal/models"
	"auth-service/internal/services"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	keycloakService *services.KeycloakService
}

func NewAuthHandler(ks *services.KeycloakService) *AuthHandler {
	return &AuthHandler{
		keycloakService: ks,
	}
}

type AuthInterface interface {
	LoginHandler(c *fiber.Ctx) error
	RegisterHandler(c *fiber.Ctx) error
	UpdateHandler(c *fiber.Ctx) error
	DeleteHandler(c *fiber.Ctx) error
	GetUserHandler(c *fiber.Ctx) error
}

func (h *AuthHandler) LoginHandler(c *fiber.Ctx) error {

	login := c.Locals("login").(models.LoginParams)

	token, err := h.keycloakService.Login(login)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "login failed",
		})
	}

	return c.JSON(token)
}

func (h *AuthHandler) RegisterHandler(c *fiber.Ctx) error {

	register := c.Locals("register").(models.RegisterParams)

	err := h.keycloakService.Register(register)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user creation failed",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "user registered successfully",
	})
}

func (h *AuthHandler) GetUserHandler(c *fiber.Ctx) error {
	
	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user ID is required",
		})
	}

	user, err := h.keycloakService.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}

func (h *AuthHandler) UpdateHandler(c *fiber.Ctx) error {
	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user ID is required",
		})
	}

	var userPayload models.UserPayload
	if err := c.BodyParser(&userPayload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	user := gocloak.User{
		ID:        gocloak.StringP(userID),
		FirstName: gocloak.StringP(userPayload.Firstname),
		LastName:  gocloak.StringP(userPayload.Lastname),
		Username:  gocloak.StringP(userPayload.Username),
		Email:     gocloak.StringP(userPayload.Email),
	}

	err := h.keycloakService.UpdateUser(userID, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "user updated successfully",
	})
}

func (h *AuthHandler) DeleteHandler(c *fiber.Ctx) error {
	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user ID is required",
		})
	}

	err := h.keycloakService.DeleteUser(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "user deleted successfully",
	})
}
