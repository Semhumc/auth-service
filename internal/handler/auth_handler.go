package handler

import (
	"auth-service/internal/models"
	"auth-service/internal/services"
	"fmt"
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
	LogoutHandler(c *fiber.Ctx) error
}

func (h *AuthHandler) LoginHandler(c *fiber.Ctx) error {
	fmt.Printf("🔐 LoginHandler called\n")
	
	// Check if login data exists in locals
	loginData := c.Locals("login")
	if loginData == nil {
		fmt.Printf("❌ No login data in locals\n")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "login data not found in request",
		})
	}
	
	login, ok := loginData.(models.LoginParams)
	if !ok {
		fmt.Printf("❌ Login data type assertion failed\n")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid login data format",
		})
	}

	fmt.Printf("📧 Login attempt for email: %s\n", login.Email)

	token, err := h.keycloakService.Login(login)
	if err != nil {
		fmt.Printf("❌ Keycloak login failed: %v\n", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "login failed",
			"details": err.Error(), // Include error details for debugging
		})
	}

	fmt.Printf("✅ Login successful!\n")

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
	fmt.Printf("👋 LogoutHandler called\n")
	
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

func (h *AuthHandler) GetProfileHandler(c *fiber.Ctx) error {
	fmt.Printf("👤 GetProfileHandler called\n")
	
	// Token'ı cookie'den al
	token := c.Cookies("access_token")
	if token == "" {
		// Header'dan da kontrol et
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			fmt.Printf("❌ No access token provided\n")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "no access token provided",
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

	fmt.Printf("🔍 Getting user profile with token\n")

	user, err := h.keycloakService.GetUserProfile(token)
	if err != nil {
		fmt.Printf("❌ Get user profile failed: %v\n", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid or expired token",
			"details": err.Error(),
		})
	}

	fmt.Printf("✅ User profile retrieved successfully\n")
	return c.JSON(user)
}

func (h *AuthHandler) RegisterHandler(c *fiber.Ctx) error {
	fmt.Printf("📝 RegisterHandler called\n")
	
	registerData := c.Locals("register")
	if registerData == nil {
		fmt.Printf("❌ No register data in locals\n")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "register data not found in request",
		})
	}
	
	register, ok := registerData.(models.RegisterParams)
	if !ok {
		fmt.Printf("❌ Register data type assertion failed\n")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid register data format",
		})
	}

	fmt.Printf("📧 Registration attempt for email: %s\n", register.Login.Email)

	err := h.keycloakService.Register(register)
	if err != nil {
		fmt.Printf("❌ Registration failed: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user creation failed",
			"details": err.Error(),
		})
	}

	fmt.Printf("✅ Registration successful!\n")
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