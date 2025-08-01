package routes

import (
	"auth-service/internal/handler"
	"auth-service/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(router fiber.Router, handler handler.AuthInterface) {

	api := router.Group("/api/v1")


	api.Post("/login", middleware.LoginMiddleware, handler.LoginHandler)
	api.Post("/register", middleware.RegisterMiddleware,handler.RegisterHandler)
	api.Get("/me", handler.GetProfileHandler) 

	user := api.Group("/user")
	user.Delete("/:id", middleware.DeleteMiddleware,handler.DeleteHandler)
	user.Put("/:id",middleware.UpdateMiddleware ,handler.UpdateHandler)
	user.Get("/:id", middleware.GetUserMiddleware,handler.GetUserHandler)

}
