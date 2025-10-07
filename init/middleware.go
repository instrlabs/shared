package initx

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupLogger(app *fiber.App) {
	app.Use(logger.New(logger.Config{
		Format:     "{\"ts\":\"${time}\",\"status\":${status},\"latency\":\"${latency}\",\"ip\":\"${ip}\",\"method\":\"${method}\",\"path\":\"${path}\",\"query\":\"${query}\",\"user_id\":\"${locals:UserID}\",\"referer\":\"${referer}\",\"user_agent\":\"${ua}\n",
		TimeFormat: "2006-01-02T15:04:05.000Z07:00",
		TimeZone:   "UTC",
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
