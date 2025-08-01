package main

import (
	"auth-service/internal/handler"
	"auth-service/internal/routes"
	"auth-service/internal/services"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Load environment variables with defaults
	keycloak_base_url := getEnvOrDefault("KEYCLOAK_BASE_URL", "http://localhost:8080")
	keycloak_realm := getEnvOrDefault("KEYCLOAK_REALM", "master")
	keycloak_client_id := getEnvOrDefault("KEYCLOAK_CLIENT_ID", "admin-cli")
	keycloak_client_secret := getEnvOrDefault("KEYCLOAK_CLIENT_SECRET", "")
	port := getEnvOrDefault("PORT", "5000")

	fmt.Printf("🚀 Starting Auth Service\n")
	fmt.Printf("   Port: %s\n", port)
	fmt.Printf("   Keycloak URL: %s\n", keycloak_base_url)
	fmt.Printf("   Keycloak Realm: %s\n", keycloak_realm)
	fmt.Printf("   Client ID: %s\n", keycloak_client_id)
	fmt.Println()

	// Create Fiber app
	app := fiber.New()

	// Create Keycloak service
	keycloakService := services.NewKeycloakService(
		keycloak_client_id,
		keycloak_client_secret,
		keycloak_realm,
		keycloak_base_url)

	// Create auth handler
	authHandler := handler.NewAuthHandler(keycloakService)

	// Setup routes
	routes.AuthRoutes(app, authHandler)

	fmt.Printf("🌐 Server starting on port %s\n", port)
	fmt.Printf("📋 Available endpoints:\n")
	fmt.Printf("   GET  http://localhost:%s/health\n", port)
	fmt.Printf("   GET  http://localhost:%s/api/v1/health\n", port)
	fmt.Printf("   POST http://localhost:%s/api/v1/login\n", port)
	fmt.Printf("   POST http://localhost:%s/api/v1/register\n", port)
	fmt.Printf("   POST http://localhost:%s/api/v1/logout\n", port)
	fmt.Printf("   GET  http://localhost:%s/api/v1/me\n", port)
	fmt.Println()
	
	// Start server
	log.Fatal(app.Listen(":" + port))
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}