package main

import (
	"auth-service/internal/handler"
	"auth-service/internal/routes"
	"auth-service/internal/services"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"

	"os"
)

var (
	keycloak_base_url      = os.Getenv("KEYCLOAK_BASE_URL")
	keycloak_realm         = os.Getenv("KEYCLOAK_REALM")
	keycloak_client_id     = os.Getenv("KEYCLOAK_CLIENT_ID")
	keycloak_client_secret = os.Getenv("KEYCLOAK_CLIENT_SECRET")
)

func main() {

	app := fiber.New()

	keycloakService := services.NewKeycloakService( 
		keycloak_client_id,
		keycloak_client_secret,
		keycloak_realm,
		keycloak_base_url)

	authHandler := handler.NewAuthHandler(keycloakService)

	
	routes.AuthRoutes(app, authHandler)

	err := app.Listen(":5000")
    if err != nil {
        log.Fatalf("Sunucu başlatılamadı: %v", err)
    }

}
