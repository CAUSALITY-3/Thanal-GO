package router

import (
	"log"
	"net/http"

	"github.com/CAUSALITY-3/Thanal-GO/router/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("thanal"))

func SetupRouter() *fiber.App {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered from panic:", r)
				c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error", "error": r})
			}
		}()
		store.Options = &sessions.Options{
			MaxAge:   60 * 24 * 60 * 60, // 60 days in seconds
			HttpOnly: true,
		}
		c.Locals("session", store)
		return c.Next()
	})

	// Static files
	app.Static("/_next/static", "./static")

	// Routes
	app.Get("/thanal", func(c *fiber.Ctx) error {
		return c.SendString("Thanal is running!!!")
	})

	// Redirect logic
	app.Use(func(c *fiber.Ctx) error {
		if !fiberPathStartsWith(c.Path(), "/thanalApi") {
			// Implement redirect logic
			return c.Redirect("/redirect-path", http.StatusMovedPermanently)
		}
		return c.Next()
	})

	routes.RegisterRoutes(app)

	// Error handling middleware
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Route not found"})
	})

	return app
}

func fiberPathStartsWith(path, prefix string) bool {
	return len(path) >= len(prefix) && path[:len(prefix)] == prefix
}
