package models

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
