package models

import (
	"sistem-manajemen-armada/internal/utils"
)

type VehicleLocation struct {
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

type GeofenceEvent struct {
	VehicleID string   `json:"vehicle_id"`
	Event     string   `json:"event"`
	Location  Location `json:"location"`
	Timestamp int64    `json:"timestamp"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (v *VehicleLocation) ValidateVehicleLocationData() []string {
	var errors []string
	if !utils.IsValidVehicleID(v.VehicleID) {
		errors = append(errors, "\"invalid vehicle_id, ex: B1234XYZ\"")
	}

	if v.Latitude < -90.0 || v.Latitude > 90.0 {
		errors = append(errors, "\"invalid latitude, must between -90.0 and 90.0\"")
	}

	if v.Longitude < -180.0 || v.Longitude > 180.0 {
		errors = append(errors, "\"invalid longitude, must between -180.0 and 180\"")
	}

	return errors
}
