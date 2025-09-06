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
	GetCurrentUserHandler(c *fiber.Ctx) error  // Yeni: Giriş yapmış kullanıcının kendi bilgilerini getirme
	UpdateCurrentUserHandler(c *fiber.Ctx) error // Yeni: Giriş yapmış kullanıcının kendi bilgilerini güncelleme
	DeleteCurrentUserHandler(c *fiber.Ctx) error // Yeni: Giriş yapmış kullanıcının kendi hesabını silme
	RefreshTokenHandler(c *fiber.Ctx) error
}

func (h *AuthHandler) LoginHandler(c *fiber.Ctx) error {
	fmt.Printf("🔐 LoginHandler called\n")
	
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

	fmt.Printf("👤 Login attempt for username: %s\n", login.Username)

	token, err := h.keycloakService.Login(login)
	if err != nil {
		fmt.Printf("❌ Keycloak login failed: %v\n", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "login failed",
			"details": err.Error(),
		})
	}

	fmt.Printf("✅ Login successful!\n")

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    token.AccessToken,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})

	return c.JSON(fiber.Map{
		"message": "login successful",
		"user":    token,
	})
}

func (h *AuthHandler) LogoutHandler(c *fiber.Ctx) error {
	fmt.Printf("👋 LogoutHandler called\n")

	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if body.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "refresh token not provided",
		})
	}

	err := h.keycloakService.Logout(body.RefreshToken)
	if err != nil {
		// Log the error but still try to clear cookies and log the user out on the client side
		fmt.Printf("⚠️ Keycloak logout failed: %v\n", err)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		HTTPOnly: true,
		MaxAge:   -1,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
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
	
	token := c.Cookies("access_token")
	if token == "" {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			fmt.Printf("❌ No access token provided\n")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "no access token provided",
			})
		}
		
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

	fmt.Printf("📧 Registration attempt for email: %s, username: %s\n", register.Email, register.Username)

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

// USER MANAGEMENT ENDPOINTS

// GET /user/:id - Belirli bir kullanıcıyı ID ile getir (Admin işlemi)
func (h *AuthHandler) GetUserHandler(c *fiber.Ctx) error {
	fmt.Printf("👥 GetUserHandler called\n")
	
	// Token kontrolü
	token := c.Locals("access_token")
	if token == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "authentication required",
		})
	}

	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user ID is required",
		})
	}

	fmt.Printf("🔍 Getting user by ID: %s\n", userID)

	user, err := h.keycloakService.GetUserByID(userID)
	if err != nil {
		fmt.Printf("❌ Get user by ID failed: %v\n", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
			"details": err.Error(),
		})
	}

	fmt.Printf("✅ User retrieved successfully\n")
	return c.JSON(user)
}

// PUT /user/:id - Belirli bir kullanıcıyı güncelle
func (h *AuthHandler) UpdateHandler(c *fiber.Ctx) error {
	fmt.Printf("✏️ UpdateHandler called\n")
	
	// Token kontrolü
	token := c.Locals("access_token")
	if token == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "authentication required",
		})
	}

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
			"details": err.Error(),
		})
	}

	fmt.Printf("🔄 Updating user ID: %s\n", userID)

	user := gocloak.User{
		ID:        gocloak.StringP(userID),
		FirstName: gocloak.StringP(userPayload.Firstname),
		LastName:  gocloak.StringP(userPayload.Lastname),
		Username:  gocloak.StringP(userPayload.Username),
		Email:     gocloak.StringP(userPayload.Email),
	}

	err := h.keycloakService.UpdateUser(userID, user)
	if err != nil {
		fmt.Printf("❌ Update user failed: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "update failed",
			"details": err.Error(),
		})
	}

	fmt.Printf("✅ User updated successfully\n")
	return c.JSON(fiber.Map{
		"message": "user updated successfully",
	})
}

// DELETE /user/:id - Belirli bir kullanıcıyı sil
func (h *AuthHandler) DeleteHandler(c *fiber.Ctx) error {
	fmt.Printf("🗑️ DeleteHandler called\n")
	
	// Token kontrolü
	token := c.Locals("access_token")
	if token == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "authentication required",
		})
	}

	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user ID is required",
		})
	}

	fmt.Printf("🗑️ Deleting user ID: %s\n", userID)

	err := h.keycloakService.DeleteUser(userID)
	if err != nil {
		fmt.Printf("❌ Delete user failed: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "delete failed",
			"details": err.Error(),
		})
	}

	fmt.Printf("✅ User deleted successfully\n")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "user deleted successfully",
	})
}

// YENİ ENDPOINT'LER - Giriş yapmış kullanıcının kendi işlemleri için

// GET /user/me - Giriş yapmış kullanıcının kendi bilgilerini getir
func (h *AuthHandler) GetCurrentUserHandler(c *fiber.Ctx) error {
	fmt.Printf("👤 GetCurrentUserHandler called\n")
	
	token := c.Locals("access_token")
	if token == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "authentication required",
		})
	}

	tokenStr, ok := token.(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token format",
		})
	}

	fmt.Printf("🔍 Getting current user profile\n")

	user, err := h.keycloakService.GetUserProfile(tokenStr)
	if err != nil {
		fmt.Printf("❌ Get current user failed: %v\n", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid or expired token",
			"details": err.Error(),
		})
	}

	fmt.Printf("✅ Current user profile retrieved successfully\n")
	return c.JSON(user)
}

// PUT /user/me - Giriş yapmış kullanıcının kendi bilgilerini güncelle
func (h *AuthHandler) UpdateCurrentUserHandler(c *fiber.Ctx) error {
	fmt.Printf("✏️ UpdateCurrentUserHandler called\n")
	
	token := c.Locals("access_token")
	if token == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "authentication required",
		})
	}

	tokenStr, ok := token.(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token format",
		})
	}

	// Önce kullanıcının kendi ID'sini al
	userProfile, err := h.keycloakService.GetUserProfile(tokenStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid or expired token",
		})
	}

	var userPayload models.UserPayload
	if err := c.BodyParser(&userPayload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
			"details": err.Error(),
		})
	}

	fmt.Printf("🔄 Updating current user\n")

	user := gocloak.User{
		ID:        userProfile.ID,
		FirstName: gocloak.StringP(userPayload.Firstname),
		LastName:  gocloak.StringP(userPayload.Lastname),
		Username:  gocloak.StringP(userPayload.Username),
		Email:     gocloak.StringP(userPayload.Email),
	}

	err = h.keycloakService.UpdateUser(*userProfile.ID, user)
	if err != nil {
		fmt.Printf("❌ Update current user failed: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "update failed",
			"details": err.Error(),
		})
	}

	fmt.Printf("✅ Current user updated successfully\n")
	return c.JSON(fiber.Map{
		"message": "profile updated successfully",
	})
}

// DELETE /user/me - Giriş yapmış kullanıcının kendi hesabını sil
func (h *AuthHandler) DeleteCurrentUserHandler(c *fiber.Ctx) error {
	fmt.Printf("🗑️ DeleteCurrentUserHandler called\n")
	
	token := c.Locals("access_token")
	if token == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "authentication required",
		})
	}

	tokenStr, ok := token.(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token format",
		})
	}

	// Önce kullanıcının kendi ID'sini al
	userProfile, err := h.keycloakService.GetUserProfile(tokenStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid or expired token",
		})
	}

	fmt.Printf("🗑️ Deleting current user account\n")

	err = h.keycloakService.DeleteUser(*userProfile.ID)
	if err != nil {
		fmt.Printf("❌ Delete current user failed: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "delete failed",
			"details": err.Error(),
		})
	}

	// Hesap silindikten sonra cookie'yi de temizle
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		HTTPOnly: true,
		MaxAge:   -1,
	})

	fmt.Printf("✅ Current user account deleted successfully\n")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "account deleted successfully",
	})
}

func (h *AuthHandler) RefreshTokenHandler(c *fiber.Ctx) error {
	fmt.Printf("🔄 RefreshTokenHandler called\n")

	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if body.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "refresh token not provided",
		})
	}


token, err := h.keycloakService.RefreshToken(body.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "failed to refresh token",
			"details": err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    token.AccessToken,
		HTTPOnly: true,
		Secure:   false, 
		SameSite: "Lax",
	})

	return c.JSON(fiber.Map{
		"message": "token refreshed successfully",
		"user":    token,
	})
}
