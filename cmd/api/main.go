package main

import (
	"log"
	"sistem-manajemen-armada/internal/config"
	"sistem-manajemen-armada/internal/repository"
	"sistem-manajemen-armada/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// enable this if you want use from local
	// godotenv.Load()
	cfg := config.LoadConfig()

	repo, err := repository.NewPostgreSQLRepository(cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPass, cfg.PostgresDB)
	if err != nil {
		log.Fatalf("Failed to init PostgreSQL repository: %v", err)
	}
	defer repo.DB.Close()

	// Init API Service
	apiService := services.NewAPIService(repo)

	router := gin.Default()
	router.GET("/vehicles/:vehicle_id/location", apiService.GetLastVehicleLocation)
	router.GET("/vehicles/:vehicle_id/history", apiService.GetVehicleLocationHistory)

	// Run the API server
	log.Printf("API service starting on port %s", cfg.APIPort)
	if err := router.Run(":" + cfg.APIPort); err != nil {
		log.Fatalf("Failed to run API server: %v", err)
	}
}
