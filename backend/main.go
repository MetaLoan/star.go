package main

import (
	"log"
	"star/api"
)

func main() {
	log.Println("🌟 Starting Star API Server...")

	router := api.SetupRouter()

	log.Println("✅ Server running on http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

