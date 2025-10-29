package initx

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

// SetupServiceSwagger sets up the Swagger endpoint with dynamic server URL
// gatewayURL is the base URL of the gateway service (e.g., "http://localhost:3000")
// servicePath is the service path (e.g., "/auth" or "/images")
func SetupServiceSwagger(app *fiber.App, gatewayURL, servicePath string) error {
	// Read the static swagger.json file
	swaggerPath := filepath.Join("static", "swagger.json")
	data, err := os.ReadFile(swaggerPath)
	if err != nil {
		return fmt.Errorf("failed to read swagger.json: %w", err)
	}

	// Parse the JSON
	var swaggerSpec map[string]interface{}
	if err := json.Unmarshal(data, &swaggerSpec); err != nil {
		return fmt.Errorf("failed to parse swagger.json: %w", err)
	}

	// Replace the servers array with dynamic server URL
	swaggerSpec["servers"] = []map[string]string{
		{
			"url":         fmt.Sprintf("%s%s", gatewayURL, servicePath),
			"description": "Gateway Service",
		},
	}

	// Convert back to JSON
	modifiedSwagger, err := json.Marshal(swaggerSpec)
	if err != nil {
		return fmt.Errorf("failed to marshal modified swagger: %w", err)
	}

	// Serve the modified swagger.json
	app.Get("/swagger", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Send(modifiedSwagger)
	})

	return nil
}

func SetupServiceHealth(app *fiber.App) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
}
