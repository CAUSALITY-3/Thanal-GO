package routes

import (
	services "github.com/CAUSALITY-3/Thanal-GO/service/user"
	"github.com/CAUSALITY-3/Thanal-GO/utils"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	userService := utils.SingletonInjector.Get("userService").(*services.UserService)
	thanalApi := app.Group("/thanalApi/users")

	thanalApi.Get("", func(c *fiber.Ctx) error {
		return userService.FindUserByEmail(c) // Ensure your service method is compatible with Fiber's context
	})
	thanalApi.Get("/getallusers", func(c *fiber.Ctx) error {
		return userService.GetAllUsers(c) // Ensure your service method is compatible with Fiber's context
	})
	thanalApi.Get("/GetUsersCache", func(c *fiber.Ctx) error {
		return userService.GetUsersCache(c) // Ensure your service method is compatible with Fiber's context
	})
	thanalApi.Post("", func(c *fiber.Ctx) error {
		return userService.CreateUser(c) // Ensure your service method is compatible with Fiber's context
	})
	thanalApi.Put("", func(c *fiber.Ctx) error {
		return userService.UpsertUser(c) // Ensure your service method is compatible with Fiber's context
	})
}
