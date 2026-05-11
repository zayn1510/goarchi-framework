package tools

import (
	"fmt"
	"math"
)

const EarthRadiusKM = 6371.0

func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := toRadians(lat2 - lat1)
	dLon := toRadians(lon2 - lon1)

	lat1Rad := toRadians(lat1)
	lat2Rad := toRadians(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return EarthRadiusKM * c
}

func toRadians(degree float64) float64 {
	return degree * math.Pi / 180
}
func FormatJarak(jarak float64) string {
	if jarak < 0.5 {
		return fmt.Sprintf("%.0f m - Terdekat", jarak*1000)
	} else if jarak < 2 {
		return fmt.Sprintf("%.2f Km - Lumayan", jarak)
	}
	return fmt.Sprintf("%.2f Km - Jauh", jarak)
}
func RoundToDecimal(value float64, places int) float64 {
	multiplier := math.Pow(10, float64(places))
	return math.Round(value*multiplier) / multiplier
}