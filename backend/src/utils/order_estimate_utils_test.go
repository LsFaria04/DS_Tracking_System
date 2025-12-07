package utils

import (
    "math"
    "testing"
)

// Allow small floating-point tolerance
func almostEqual(a, b, tolerance float64) bool {
    return math.Abs(a-b) <= tolerance
}

func TestHaversine_SamePoint(t *testing.T) {
    // Distance between identical coordinates should be 0
    d := haversine(41.14961, -8.61099, 41.14961, -8.61099)
    if !almostEqual(d, 0, 1e-9) {
        t.Errorf("expected 0, got %f", d)
    }
}

func TestHaversine_KnownDistance(t *testing.T) {
    // Porto (41.14961, -8.61099) to Lisbon (38.7169, -9.1390)
    // Known distance ≈ 274 km
    d := haversine(41.14961, -8.61099, 38.7169, -9.1390)
    if !almostEqual(d, 274, 1.0) {
        t.Errorf("expected ~274 km, got %f", d)
    }
}

func TestEstimateDeliveryTime(t *testing.T) {
    // Porto to Lisbon, average speed 100 km/h
    // Distance ≈ 274 km → travel time ≈ 2.74 h + 12 h buffer = ~14.74 h
    et := EstimateDeliveryTime(41.14961, -8.61099, 38.7169, -9.1390, 100)
    if !almostEqual(et, 14.74, 0.5) {
        t.Errorf("expected ~14.74 h, got %f", et)
    }
}

func TestEstimateDeliveryTime_ZeroDistance(t *testing.T) {
    // Same point, should return only the 12h buffer
    et := EstimateDeliveryTime(41.14961, -8.61099, 41.14961, -8.61099, 100)
    if !almostEqual(et, 12, 1e-9) {
        t.Errorf("expected 12, got %f", et)
    }
}

func TestEstimateDeliveryTime_DifferentSpeeds(t *testing.T) {
    // Porto to Lisbon, slower speed
    etSlow := EstimateDeliveryTime(41.14961, -8.61099, 38.7169, -9.1390, 50)
    etFast := EstimateDeliveryTime(41.14961, -8.61099, 38.7169, -9.1390, 200)

    if etSlow <= etFast {
        t.Errorf("expected slower speed to give longer time, got slow=%f fast=%f", etSlow, etFast)
    }
}
