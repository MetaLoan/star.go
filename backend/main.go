package main

import (
	"log"
	"star/api"
	"star/astro"
)

func main() {
	log.Println("🌟 Starting Star API Server...")

	// 验证 Swiss Ephemeris 是否可用（唯一数据源）
	if err := astro.ValidateSwissEphemeris(); err != nil {
		log.Fatalf("❌ %v", err)
	}
	log.Printf("✅ Data Source: %s", astro.GetDataSource())

	// 确保在程序结束时关闭 Swiss Ephemeris
	defer astro.CloseSwissEphemeris()

	router := api.SetupRouter()

	log.Println("✅ Server running on http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

