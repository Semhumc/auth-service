package routes

import (
	"auth-service/internal/handler"
	"auth-service/internal/middleware"
	"auth-service/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func AuthRoutes(app *fiber.App, handler handler.AuthInterface, keycloakService *services.KeycloakService) {
	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",
		AllowCredentials: true,
		ExposeHeaders:    "Set-Cookie",
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "auth-service",
			"version": "1.0.0",
		})
	})

	// API routes
	api := app.Group("/api/v1")

	// Health check for API
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "auth-service API",
			"version": "1.0.0",
		})
	})

	authTokenMiddleware := middleware.NewAuthTokenMiddleware(keycloakService)

	// AUTH ENDPOINTS (Token gerektirmeyen)
	api.Post("/login", middleware.LoginMiddleware, handler.LoginHandler)
	api.Post("/register", middleware.RegisterMiddleware, handler.RegisterHandler)
	api.Post("/logout", handler.LogoutHandler)
	api.Post("/refresh", handler.RefreshTokenHandler)
	api.Get("/me", handler.GetProfileHandler) // Eski endpoint, uyumluluk için

	// USER MANAGEMENT ENDPOINTS (Token gerektiren)
	user := api.Group("/user")
	
	// Giriş yapmış kullanıcının kendi işlemleri (Token ile)
	user.Get("/me", authTokenMiddleware, handler.GetCurrentUserHandler)
	user.Put("/me", authTokenMiddleware, handler.UpdateCurrentUserHandler)
	user.Delete("/me", authTokenMiddleware, handler.DeleteCurrentUserHandler)
	
	// Admin seviyesi işlemler (ID ile) - Token gerekli
	user.Get("/:id", authTokenMiddleware, middleware.GetUserMiddleware, handler.GetUserHandler)
	user.Put("/:id", authTokenMiddleware, middleware.UpdateMiddleware, handler.UpdateHandler)
	user.Delete("/:id", authTokenMiddleware, middleware.DeleteMiddleware, handler.DeleteHandler)

	// Debug endpoint
	api.Get("/test-cors", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "CORS is working!",
			"origin":  c.Get("Origin"),
			"method":  c.Method(),
		})
	})

	// Catch-all route
	app.Use("*", func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"error":  "Route not found",
			"method": c.Method(),
			"path":   c.Path(),
		})
	})
}
