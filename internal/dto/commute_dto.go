package dto

import "github.com/google/uuid"

type CreateCommuteRequest struct {
	DeviceID    string  `json:"device_id" binding:"required"`
	Name        string  `json:"name"`
	HomeLng     float64 `json:"home_lng" binding:"required"`
	HomeLat     float64 `json:"home_lat" binding:"required"`
	OfficeLng   float64 `json:"office_lng" binding:"required"`
	OfficeLat   float64 `json:"office_lat" binding:"required"`
	Vehicle     string  `json:"vehicle" binding:"required,oneof=car motorcycle"`
	FuelPrice   int32   `json:"fuel_price" binding:"required,gt=0"`
	DaysPerWeek int16   `json:"days_per_week" binding:"required,min=1,max=7"`
}

type UpdateCommuteRequest struct {
	Name string `json:"name"`
}

type Commute struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	HomeLng        float64   `json:"home_lng"`
	HomeLat        float64   `json:"home_lat"`
	OfficeLng      float64   `json:"office_lng"`
	OfficeLat      float64   `json:"office_lat"`
	DistanceKm     float64   `json:"distance_km"`
	DurationMin    float64   `json:"duration_min"`
	Vehicle        string    `json:"vehicle"`
	FuelPrice      int32     `json:"fuel_price"`
	DaysPerWeek    int16     `json:"days_per_week"`
	AnnualCostRp   int64     `json:"annual_cost_rp"`
	AnnualMinutes  int64     `json:"annual_minutes"`
	AnnualHours    float64   `json:"annual_hours"`
	AnnualWorkdays float64   `json:"annual_workdays"`
	CreatedAt      string    `json:"created_at"`
}

type ListCommutesResponse struct {
	Commutes []Commute `json:"commutes"`
	Total    int       `json:"total"`
}
