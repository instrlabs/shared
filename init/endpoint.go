package initx

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// SetupServiceSwagger sets up the Swagger endpoint with dynamic server URL
// gatewayURL is the base URL of the gateway service (e.g., "http://localhost:3000")
// servicePath is the service path (e.g., "/auth" or "/images")
// Note: Errors are logged internally and do not need to be handled by the caller
func SetupServiceSwagger(app *fiber.App, gatewayURL, servicePath string) {
	// Read the static swagger.json file
	swaggerPath := filepath.Join("static", "swagger.json")
	data, err := os.ReadFile(swaggerPath)
	if err != nil {
		log.Warnf("Failed to read swagger.json: %v", err)
		return
	}

	// Parse the JSON
	var swaggerSpec map[string]interface{}
	if err := json.Unmarshal(data, &swaggerSpec); err != nil {
		log.Warnf("Failed to parse swagger.json: %v", err)
		return
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
		log.Warnf("Failed to marshal modified swagger: %v", err)
		return
	}

	// Serve the modified swagger.json
	app.Get("/swagger", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Send(modifiedSwagger)
	})
}

func SetupServiceHealth(app *fiber.App) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
}
