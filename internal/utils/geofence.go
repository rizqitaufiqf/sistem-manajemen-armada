package utils

import (
	"log"
	"math"
)

const earthRadius = 6371000

// IsInGeofence checks if a point is within a radius (in meters) of a center point
func IsInGeofence(lat1, lon1, lat2, lon2, radius float64) bool {
	// Convert degrees to radians
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// Haversine formula
	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad
	a := math.Pow(math.Sin(dLat/2), 2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Pow(math.Sin(dLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadius * c

	log.Printf("Distance => %.2f meters", distance)

	return distance <= radius
}
