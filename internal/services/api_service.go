package services

import (
	"net/http"
	"sistem-manajemen-armada/internal/repository"
	"strconv"

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

func (s *APIService) GetVehicleLocationHistory(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")
	if vehicleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle_id is required"})
		return
	}

	timeStart := c.Query("start")
	timeEnd := c.Query("end")

	if timeStart == "" || timeEnd == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start and end timestamps are required"})
		return
	}

	start, err := strconv.ParseInt(timeStart, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start timestamp"})
		return
	}
	end, err := strconv.ParseInt(timeEnd, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end timestamp"})
		return
	}

	locations, err := s.repo.GetVehicleLocationHistory(vehicleID, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	if len(locations) == 0 {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	c.JSON(http.StatusOK, locations)
}
