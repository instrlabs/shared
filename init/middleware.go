package initx

import (
	"os"
	"strings"
	"time"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupLogger(app *fiber.App) {
	app.Use(logger.New(logger.Config{
		Format:     `{"time":"${time}","host":"${header:x-user-host}","ip":"${header:x-user-ip}","method":"${method}","path":"${path}","status":${status},"latency":"${latency}","userAgent":"${header:x-user-agent}"}` + "\n",
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

		userId := c.Get("x-user-id")
		if userId != "" {
			c.Locals("userId", userId)
		}

		// Extract and set device info from headers for session binding
		ipAddress := c.Get("x-user-ip")
		if ipAddress != "" {
			c.Locals("userIP", ipAddress)
		}

		userAgent := c.Get("x-user-agent")
		if userAgent != "" {
			c.Locals("userAgent", userAgent)
		}

		if !isPublic(c.Path()) && userId == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
				"errors":  nil,
				"data":    nil,
			})
		}

		return c.Next()
	})
}

func SetupPrometheus(app *fiber.App) {
	appName := GetEnv("APP_NAME", "-")
	prom := fiberprometheus.New(appName)
	prom.RegisterAt(app, "/metrics")
	app.Use(prom.Middleware)
}
