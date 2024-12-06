package routes

import (
	services "github.com/CAUSALITY-3/Thanal-GO/service/user"
	"github.com/CAUSALITY-3/Thanal-GO/utils"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(thanalApiRoute fiber.Router) {

	userRoute := thanalApiRoute.Group("/users")

	userService := utils.SingletonInjector.Get("userService").(*services.UserService)

	userRoute.Get("", func(c *fiber.Ctx) error {
		return userService.FindUserByEmail(c)
	})
	userRoute.Get("/getallusers", func(c *fiber.Ctx) error {
		return userService.GetAllUsers(c)
	})
	userRoute.Get("/GetUsersCache", func(c *fiber.Ctx) error {
		return userService.GetUsersCache(c)
	})
	userRoute.Post("", func(c *fiber.Ctx) error {
		return userService.CreateUser(c)
	})
	userRoute.Put("", func(c *fiber.Ctx) error {
		return userService.UpsertUser(c)
	})
	userRoute.Put("/UpdateUserOrder", func(c *fiber.Ctx) error {
		return userService.UpdateUserOrder(c)
	})
	userRoute.Put("/addToBag", func(c *fiber.Ctx) error {
		return userService.AddToBag(c)
	})
	userRoute.Put("/removeFromBag", func(c *fiber.Ctx) error {
		return userService.RemoveFromBag(c)
	})
	userRoute.Put("/favoriteItem", func(c *fiber.Ctx) error {
		return userService.FavoriteItem(c)
	})
	userRoute.Put("/unFavoriteItem", func(c *fiber.Ctx) error {
		return userService.UnfavoriteItem(c)
	})
}
