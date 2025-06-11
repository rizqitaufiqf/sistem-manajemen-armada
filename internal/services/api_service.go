package services

import (
	"net/http"
	"sistem-manajemen-armada/internal/repository"

	"github.com/gin-gonic/gin"
)

type APIService struct {
	repo *repository.PostgreSQLRepository
}

func NewAPIService(repo *repository.PostgreSQLRepository) *APIService {
	return &APIService{repo: repo}
}

func (s *APIService) GetLastVehicleLocation(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")
	if vehicleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle_id is required"})
		return
	}

	vehicleLocation, err := s.repo.GetLastVehicleLocation(vehicleID)
	if err != nil {
		if err.Error() == "no location found for vehicle "+vehicleID {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vehicleLocation)
}
