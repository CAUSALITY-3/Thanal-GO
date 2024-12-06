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
		return userService.FindUserByEmail(c)
	})
	thanalApi.Get("/getallusers", func(c *fiber.Ctx) error {
		return userService.GetAllUsers(c)
	})
	thanalApi.Get("/GetUsersCache", func(c *fiber.Ctx) error {
		return userService.GetUsersCache(c)
	})
	thanalApi.Post("", func(c *fiber.Ctx) error {
		return userService.CreateUser(c)
	})
	thanalApi.Put("", func(c *fiber.Ctx) error {
		return userService.UpsertUser(c)
	})
	thanalApi.Put("/UpdateUserOrder", func(c *fiber.Ctx) error {
		return userService.UpdateUserOrder(c)
	})
	thanalApi.Put("/addToBag", func(c *fiber.Ctx) error {
		return userService.AddToBag(c)
	})
	thanalApi.Put("/removeFromBag", func(c *fiber.Ctx) error {
		return userService.RemoveFromBag(c)
	})
	thanalApi.Put("/favoriteItem", func(c *fiber.Ctx) error {
		return userService.FavoriteItem(c)
	})
	thanalApi.Put("/unFavoriteItem", func(c *fiber.Ctx) error {
		return userService.UnfavoriteItem(c)
	})
}
