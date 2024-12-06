package routes

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	thanalApiRoute := app.Group("/thanalApi")
	UserRoutes(thanalApiRoute)
}
