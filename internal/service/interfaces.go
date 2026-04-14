package service

import (
	"context"
)

// Coord represents a geographic coordinate.
type Coord struct {
	Lat, Lng float64
}

// RouteResult is a provider-agnostic route calculation result.
type RouteResult struct {
	DistanceKm  float64
	DurationMin float64
	Coordinates [][]float64 // [lng, lat] pairs
}

// RoutingService abstracts external route calculation providers.
type RoutingService interface {
	GetRoute(ctx context.Context, profile string, start, end Coord) (*RouteResult, error)
}
