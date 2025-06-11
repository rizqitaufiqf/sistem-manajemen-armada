package utils

import "regexp"

func IsValidVehicleID(vehicleID string) bool {
	pattern := `^[a-zA-Z]{1,2}\d{1,4}[a-zA-Z]{1,3}$`
	regex := regexp.MustCompile(pattern)

	return regex.MatchString(vehicleID)
}
