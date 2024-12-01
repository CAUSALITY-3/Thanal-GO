package routes

import (
	services "github.com/CAUSALITY-3/Thanal-GO/service/user"
	"github.com/CAUSALITY-3/Thanal-GO/utils"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	userService := utils.SingletonInjector.Get("userService").(*services.UserService)
	thanalApi := r.Group("/thanalApi")

	{
		thanalApi.GET("/users", userService.FindUserByEmail)
		thanalApi.POST("/users", userService.CreateUser)
	}

	// {
	//     // Product routes
	//     thanalApi.GET("/products", userService.FindUserByEmail)
	//     thanalApi.POST("/products", createProductHandler)
	//     thanalApi.PUT("/products/:id", updateProductHandler)
	//     thanalApi.DELETE("/products/:id", deleteProductHandler)

	//     // Features routes
	//     thanalApi.GET("/features", getFeaturesHandler)
	//     thanalApi.POST("/features", createFeatureHandler)
	//     thanalApi.PUT("/features/:id", updateFeatureHandler)
	//     thanalApi.DELETE("/features/:id", deleteFeatureHandler)

	//     // Auth routes
	//     thanalApi.POST("/auth", loginHandler)
	//     thanalApi.POST("/auth/register", registerHandler)

	//     // User routes with authentication middleware
	//     thanalApi.Use(authenticateMiddleware)
	//     {
	//         thanalApi.GET("/users", getUsersHandler)
	//         thanalApi.POST("/users", createUserHandler)
	//         thanalApi.PUT("/users/:id", updateUserHandler)
	//         thanalApi.DELETE("/users/:id", deleteUserHandler)

	//         thanalApi.GET("/payments", getPaymentsHandler)
	//         thanalApi.POST("/payments", createPaymentHandler)
	//         thanalApi.PUT("/payments/:id", updatePaymentHandler)
	//         thanalApi.DELETE("/payments/:id", deletePaymentHandler)

	//         thanalApi.GET("/orders", getOrdersHandler)
	//         thanalApi.POST("/orders", createOrderHandler)
	//         thanalApi.PUT("/orders/:id", updateOrderHandler)
	//         thanalApi.DELETE("/orders/:id", deleteOrderHandler)
	//     }

	//     // Other routes
	//     thanalApi.GET("/images", getImagesHandler)
	//     thanalApi.POST("/upload", uploadHandler)
	// }
}
