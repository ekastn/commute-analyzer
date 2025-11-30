package dto

import "github.com/paulmach/orb"

type CreateCommuteRequest struct {
	DeviceID    string  `json:"device_id" binding:"required"`
	Name        string  `json:"name"`
	HomeLat     float64 `json:"home_lat" binding:"required"`
	HomeLng     float64 `json:"home_lng" binding:"required"`
	OfficeLat   float64 `json:"office_lat" binding:"required"`
	OfficeLng   float64 `json:"office_lng" binding:"required"`
	Vehicle     string  `json:"vehicle" binding:"required"`
	FuelPrice   int     `json:"fuel_price" binding:"required"`
	DaysPerWeek int     `json:"days_per_week" binding:"required"`
}

type UpdateCommuteRequest struct {
	Name string `json:"name" binding:"required"`
}

type Commute struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	HomeLat        float64         `json:"home_lat"`
	HomeLng        float64         `json:"home_lng"`
	OfficeLat      float64         `json:"office_lat"`
	OfficeLng      float64         `json:"office_lng"`
	RouteGeometry  *orb.LineString `json:"route_geometry,omitempty"`
	DistanceKm     float64         `json:"distance_km"`
	DurationMin    float64         `json:"duration_min"`
	Vehicle        string          `json:"vehicle"`
	FuelPrice      int32           `json:"fuel_price"`
	DaysPerWeek    int32           `json:"days_per_week"`
	AnnualCostRp   int64           `json:"annual_cost_rp"`
	AnnualMinutes  int64           `json:"annual_minutes"`
	AnnualHours    float64         `json:"annual_hours"`
	AnnualWorkdays float64         `json:"annual_workdays"`
	CreatedAt      string          `json:"created_at"`
}

type ListCommutesResponse struct {
	Commutes []Commute `json:"commutes"`
	Total    int       `json:"total"`
}
