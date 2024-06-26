package utils

import (
	"beli_mang/db/entities"
	"beli_mang/responses"
	"math"
)

const earthRadius = 6371   // Earth's radius in kilometers
const deliverySpeed = 40.0 // Speed in km/h
// Define the area as a constant
const area = 3.0 // km²

// Calculate the diameter for a circle with an area of 3 km²
var diameter = 2 * math.Sqrt(area/math.Pi)

// Haversine formula to calculate distance between two points
func haversine(lat1, lon1, lat2, lon2 float64) (float64, error) {
	lat1Rad := lat1 * (math.Pi / 180)
	lon1Rad := lon1 * (math.Pi / 180)
	lat2Rad := lat2 * (math.Pi / 180)
	lon2Rad := lon2 * (math.Pi / 180)

	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad

	a := math.Pow(math.Sin(dLat/2), 2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Pow(math.Sin(dLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadius * c

	return distance, nil
}

func IsWithin3Km2(points []entities.RoutePoint) (bool, error) {
	if len(points) == 0 {
		return false, responses.NewInternalServerError("no points provided")
	}

	// Ending point is the last element
	end := points[len(points)-1]

	// Calculate the distance between each point and the ending point
	for _, point := range points {
		distance, err := haversine(point.Latitude, point.Longitude, end.Latitude, end.Longitude)
		if err != nil {
			return false, responses.NewInternalServerError("Error calculating distance:" + err.Error())
		}
		if distance > diameter {
			return false, nil
		}
	}

	return true, nil
}

func NearestNeighborTSP(locations []entities.RoutePoint) ([]int, float64, error) {
	n := len(locations)

	// Check for edge cases
	if n < 2 {
		return nil, 0, responses.NewInternalServerError("not enough locations to form a route")
	}

	start := 0
	end := n - 1

	visited := make([]bool, n)
	route := make([]int, 0)
	totalDistance := 0.0

	current := start
	route = append(route, current)
	visited[current] = true

	for len(route) < n-1 {
		nearest := -1
		minDistance := math.MaxFloat64

		for i := range locations {
			if !visited[i] && i != end {
				dist, err := haversine(locations[current].Latitude, locations[current].Longitude, locations[i].Latitude, locations[i].Longitude)
				if err != nil {
					return nil, 0, err
				}
				if dist < minDistance {
					minDistance = dist
					nearest = i
				}
			}
		}

		if nearest != -1 {
			route = append(route, nearest)
			visited[nearest] = true
			totalDistance += minDistance
			current = nearest
		}
	}

	// Add the end location
	route = append(route, end)
	finalLegDistance, err := haversine(locations[current].Latitude, locations[current].Longitude, locations[end].Latitude, locations[end].Longitude)
	if err != nil {
		return nil, 0, err
	}
	totalDistance += finalLegDistance

	return route, totalDistance, nil
}

func EstimatedDeliveryTimeInMinutes(totalDistance float64) float64 {
	// Speed is given in km/h, converting to minutes
	timeInHours := totalDistance / deliverySpeed
	timeInMinutes := timeInHours * 60
	return timeInMinutes
}
