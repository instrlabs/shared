package initx

import (
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupLogger(app *fiber.App) {
	app.Use(logger.New(logger.Config{
		Format:     `{"time":"${time}","host":"${header:X-Forwarded-Host}","ip":"${header:X-Forwarded-For}","method":"${method}","path":"${path}","status":${status},"latency":"${latency}","userAgent":"${header:X-User-Agent}"}` + "\n",
		TimeFormat: time.RFC3339Nano,
		TimeZone:   "UTC",
		Output:     os.Stdout,
	}))
}

func SetupAuthenticated(app *fiber.App, whitelist []string) {
	app.Use(func(c *fiber.Ctx) error {
		isPublic := func(path string) bool {
			for _, prefix := range whitelist {
				if path == prefix || strings.HasPrefix(path, prefix) {
					return true
				}
			}
			return false
		}

		isAuthenticated := c.Get("X-Authenticated") == "true"
		if isAuthenticated {
			userId := c.Get("X-User-Id")
			c.Locals("UserID", userId)
			roles := c.Get("X-User-Roles")
			c.Locals("Roles", roles)
		}

		if !isPublic(c.Path()) && !isAuthenticated {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
				"errors":  nil,
				"data":    nil,
			})
		}

		return c.Next()
	})
}
