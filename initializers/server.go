package initializers

import (
	"os"

	"github.com/CAUSALITY-3/Thanal-GO/router"
)

func ServerInitialize() {
	router := router.SetupRouter()
	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	router.Run(":" + port)
}
