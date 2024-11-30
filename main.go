package main

import (
	"log"

	"github.com/CAUSALITY-3/Thanal-GO/initializers"
	database "github.com/CAUSALITY-3/Thanal-GO/models/DB"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load() // Load .env file
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	err = database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer database.DisconnectDB()
	initializers.InjectServices()
	initializers.ServerInitialize()
}
