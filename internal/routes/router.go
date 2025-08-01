package routes

import (
	"auth-service/internal/handler"
	"auth-service/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func AuthRoutes(app *fiber.App, handler handler.AuthInterface) {
	// Add request logging
	app.Use(logger.New())

	// CORS middleware - configured for React development
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

	// Auth endpoints
	api.Post("/login", middleware.LoginMiddleware, handler.LoginHandler)
	api.Post("/register", middleware.RegisterMiddleware, handler.RegisterHandler)
	//api.Post("/logout", handler.LogoutHandler)
	api.Get("/me", handler.GetProfileHandler)

	// User management endpoints
	user := api.Group("/user")
	user.Get("/:id", middleware.GetUserMiddleware, handler.GetUserHandler)
	user.Put("/:id", middleware.UpdateMiddleware, handler.UpdateHandler)
	user.Delete("/:id", middleware.DeleteMiddleware, handler.DeleteHandler)

	// Debug endpoint to test CORS
	api.Get("/test-cors", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "CORS is working!",
			"origin":  c.Get("Origin"),
			"method":  c.Method(),
		})
	})

	// Catch-all route for debugging
	app.Use("*", func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"error":  "Route not found",
			"method": c.Method(),
			"path":   c.Path(),
		})
	})
}