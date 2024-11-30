package initializers

import (
	"log"
	"net/http"
	"os"

	services "github.com/CAUSALITY-3/Thanal-GO/service/user"
	"github.com/CAUSALITY-3/Thanal-GO/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/rs/cors"
)

var store = sessions.NewCookieStore([]byte("thanal"))

func ServerInitialize() {
	log.Println("initializeServer")
	router := gin.Default()

	// Middleware setup
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(func(c *gin.Context) {
		store.Options = &sessions.Options{
			MaxAge:   60 * 24 * 60 * 60, // 60 days in seconds
			HttpOnly: true,
		}
		c.Set("session", store)
		c.Next()
	})
	corsMiddleware := cors.Default()
	router.Use(func(c *gin.Context) {
		corsMiddleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Next()
		})).ServeHTTP(c.Writer, c.Request)
	})
	// Static files
	router.Static("/_next/static", "/static")

	// Routes
	router.GET("/thanal", func(c *gin.Context) {
		c.String(http.StatusOK, "Thanal is running!!!")
	})

	router.GET("/writeCache", func(c *gin.Context) {
		// Implement writeCacheToFile logic
		status := "Cache written to file"
		c.String(http.StatusOK, status)
	})

	router.GET("/getRedirectCache", func(c *gin.Context) {
		// Implement cache retrieval logic
		keys := []string{"key1", "key2", "key3"} // Replace with actual cache keys
		c.JSON(http.StatusOK, keys)
	})

	router.GET("/generateAndLoadCache", func(c *gin.Context) {
		// Implement generateAndLoadCache logic
		status := "Cache generated and loaded"
		c.String(http.StatusOK, status)
	})

	router.GET("/loadCache", func(c *gin.Context) {
		// Implement loadCache logic
		c.String(http.StatusOK, "Loaded")
	})

	router.GET("/getUsersCache", func(c *gin.Context) {
		// Implement usersCache logic
		data := map[string]string{"user1": "data1", "user2": "data2"} // Replace with actual user cache data
		c.JSON(http.StatusOK, data)
	})

	// Redirect logic
	router.Use(func(c *gin.Context) {
		if !ginPathStartsWith(c.Request.URL.Path, "/thanalApi") {
			// Implement redirect logic
			c.Redirect(http.StatusMovedPermanently, "/redirect-path")
			return
		}
		c.Next()
	})

	// Subroutes
	api := router.Group("/thanalApi/user")
	userService := utils.SingletonInjector.Get("userService").(*services.UserService)
	{
		api.GET("", userService.FindUserByEmail)
		api.POST("/features", userService.FindUserByEmail)

	}

	// Error handling middleware
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	router.Run(":" + port)
}

func ginPathStartsWith(path, prefix string) bool {
	return len(path) >= len(prefix) && path[:len(prefix)] == prefix
}
