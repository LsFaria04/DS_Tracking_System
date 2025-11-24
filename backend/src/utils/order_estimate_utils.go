package utils

import (
    "math"
)

const EarthRadiusKm = 6371.0

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
    dLat := (lat2 - lat1) * math.Pi / 180.0
    dLon := (lon2 - lon1) * math.Pi / 180.0

    a := math.Sin(dLat/2)*math.Sin(dLat/2) +
        math.Cos(lat1*math.Pi/180.0)*math.Cos(lat2*math.Pi/180.0) *
            math.Sin(dLon/2)*math.Sin(dLon/2)

    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    return EarthRadiusKm * c
}

func EstimateDeliveryTime(lat1, lon1, lat2, lon2, avgSpeedKmH float64) float64 {
	distanceKm := haversine(lat1,lon1, lat2, lon2)
    return distanceKm / avgSpeedKmH + 12 //adds an extra time to predict possible time in storages
}
