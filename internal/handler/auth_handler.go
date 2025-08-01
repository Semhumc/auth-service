package handler

import (
	"auth-service/internal/models"
	"auth-service/internal/services"
	"strings"

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
	GetProfileHandler(c *fiber.Ctx) error
}

func (h *AuthHandler) LoginHandler(c *fiber.Ctx) error {

	login := c.Locals("login").(models.LoginParams)

	token, err := h.keycloakService.Login(login)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "login failed",
		})
	}

	// Token'ı HTTP-only cookie olarak set et (güvenlik için)
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    token.AccessToken,
		HTTPOnly: true,
		Secure:   false, // Development için false, production'da true olmalı
		SameSite: "Lax",
	})

	return c.JSON(fiber.Map{
		"message": "login successful",
		"user":    token,
	})
}

func (h *AuthHandler) LogoutHandler(c *fiber.Ctx) error {
	// Cookie'yi temizle
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		HTTPOnly: true,
		MaxAge:   -1,
	})

	return c.JSON(fiber.Map{
		"message": "logout successful",
	})
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
func (h *AuthHandler) GetProfileHandler(c *fiber.Ctx) error {
	// Token'ı cookie'den al
	token := c.Cookies("access_token")
	if token == "" {
		// Header'dan da kontrol et
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "no access token provided",
			})
		}
		
		// Bearer token formatını kontrol et
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization header format",
			})
		}
		token = parts[1]
	}

	user, err := h.keycloakService.GetUserProfile(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid or expired token",
		})
	}

	return c.JSON(user)
}