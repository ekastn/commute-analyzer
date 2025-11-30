package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ekastn/commute-analyzer/internal/dto"
	"github.com/ekastn/commute-analyzer/internal/store"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkb"
)

type CommuteService struct {
	store       store.Querier
	ors         ORSClient
	userService *UserService
}

func NewCommuteService(store store.Querier, ors ORSClient, userService *UserService) *CommuteService {
	return &CommuteService{store: store, ors: ors, userService: userService}
}

func (s *CommuteService) CreateCommute(ctx context.Context, req dto.CreateCommuteRequest) (*dto.Commute, error) {
	userID, err := s.userService.GetOrCreateUser(ctx, req.DeviceID)
	if err != nil {
		return nil, err
	}

	// Map vehicle to ORS profile
	profile := "driving-car" // Default profile

	switch req.Vehicle {
	case "car":
		profile = "driving-car"
	case "motorcycle":
		profile = "cycling-regular"
	// case "bicycle": // Example for future expansion
	// 	profile = "cycling-regular"
	}

	route, err := s.ors.FetchRoute(ctx, profile, req.HomeLng, req.HomeLat, req.OfficeLng, req.OfficeLat)
	if err != nil {
		return nil, fmt.Errorf("fetch route: %w", err)
	}

	summary := route.Features[0].Properties.Summary
	distanceKm := summary.Distance / 1000
	durationMin := summary.Duration / 60

	lineString := coordsToLineString(route.Features[0].Geometry.Coordinates)

	efficiency := map[string]float64{"car": 10.0, "motorcycle": 2.5}[req.Vehicle]
	roundTrip := distanceKm * 2
	dailyCost := (roundTrip * efficiency / 100) * float64(req.FuelPrice)
	annualCost := int64(dailyCost * float64(req.DaysPerWeek) * 52.142857)
	annualMinutes := int64(durationMin*2) * int64(req.DaysPerWeek) * 52

	name := req.Name

	homePoint := orb.Point{req.HomeLng, req.HomeLat}
	officePoint := orb.Point{req.OfficeLng, req.OfficeLat}

	homeWKB, _ := wkb.Marshal(homePoint)
	officeWKB, _ := wkb.Marshal(officePoint)
	routeWKB, _ := wkb.Marshal(lineString)

	row, err := s.store.CreateCommute(ctx, store.CreateCommuteParams{
		UserID:        userID,
		Name:          &name,
		HomePoint:     homeWKB,
		OfficePoint:   officeWKB,
		RouteGeometry: routeWKB,
		DistanceKm:    distanceKm,
		DurationMin:   durationMin,
		Vehicle:       req.Vehicle,
		FuelPrice:     int32(req.FuelPrice),
		DaysPerWeek:   int16(req.DaysPerWeek),
		AnnualCost:    annualCost,
		AnnualMinutes: annualMinutes,
	})
	if err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}

	return &dto.Commute{
		ID:             row.ID.String(),
		Name:           name,
		HomeLng:        req.HomeLng,
		HomeLat:        req.HomeLat,
		OfficeLng:      req.OfficeLng,
		OfficeLat:      req.OfficeLat,
		RouteGeometry:  &lineString,
		DistanceKm:     distanceKm,
		DurationMin:    durationMin,
		Vehicle:        req.Vehicle,
		FuelPrice:      int32(req.FuelPrice),
		DaysPerWeek:    int32(req.DaysPerWeek),
		AnnualCostRp:   annualCost,
		AnnualMinutes:  annualMinutes,
		AnnualHours:    float64(annualMinutes) / 60,
		AnnualWorkdays: float64(annualMinutes) / (60 * 8),
		CreatedAt:      row.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *CommuteService) ListCommutes(ctx context.Context, deviceID string) (*dto.ListCommutesResponse, error) {
	userID, err := s.userService.GetOrCreateUser(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	rows, err := s.store.ListCommutesByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	commutes := make([]dto.Commute, len(rows))
	for i, r := range rows {
		var name string

		if r.Name != nil {
			name = *r.Name
		}

		// Unmarshal WKB
		homeGeom, _ := wkb.Unmarshal(r.HomePoint.([]byte))
		officeGeom, _ := wkb.Unmarshal(r.OfficePoint.([]byte))
		homePoint := homeGeom.(orb.Point)
		officePoint := officeGeom.(orb.Point)

		routeGeom, _ := wkb.Unmarshal(r.RouteGeometry.([]byte))
		routeLineString, _ := routeGeom.(orb.LineString)

		commutes[i] = dto.Commute{
			ID:             r.ID.String(),
			Name:           name,
			HomeLng:        homePoint.Lon(),
			HomeLat:        homePoint.Lat(),
			OfficeLng:      officePoint.Lon(),
			OfficeLat:      officePoint.Lat(),
			RouteGeometry:  &routeLineString,
			DistanceKm:     r.DistanceKm,
			DurationMin:    r.DurationMin,
			Vehicle:        r.Vehicle,
			FuelPrice:      r.FuelPrice,
			DaysPerWeek:    int32(r.DaysPerWeek),
			AnnualCostRp:   r.AnnualCost,
			AnnualMinutes:  r.AnnualMinutes,
			AnnualHours:    float64(r.AnnualMinutes) / 60,
			AnnualWorkdays: float64(r.AnnualMinutes) / (60 * 8),
			CreatedAt:      r.CreatedAt.Format(time.RFC3339),
		}
	}

	return &dto.ListCommutesResponse{Commutes: commutes, Total: len(commutes)}, nil
}

func (s *CommuteService) UpdateCommute(ctx context.Context, id uuid.UUID, req dto.UpdateCommuteRequest) (*dto.Commute, error) {
	id, err := s.store.UpdateCommute(ctx, store.UpdateCommuteParams{
		ID:   id,
		Name: &req.Name,
	})
	if err != nil {
		return nil, err
	}

	row, err := s.store.GetCommute(ctx, id)
	if err != nil {
		return nil, err
	}

	hGeom, _ := wkb.Unmarshal(row.HomePoint.([]byte))
	oGeom, _ := wkb.Unmarshal(row.OfficePoint.([]byte))
	hPoint := hGeom.(orb.Point)
	oPoint := oGeom.(orb.Point)

	rGeom, _ := wkb.Unmarshal(row.RouteGeometry.([]byte))
	rLineString, _ := rGeom.(orb.LineString)

	commute := &dto.Commute{
		ID:             row.ID.String(),
		Name:           *row.Name,
		HomeLng:        hPoint.Lon(),
		HomeLat:        hPoint.Lat(),
		OfficeLng:      oPoint.Lon(),
		OfficeLat:      oPoint.Lat(),
		RouteGeometry:  &rLineString,
		DistanceKm:     row.DistanceKm,
		DurationMin:    row.DurationMin,
		Vehicle:        row.Vehicle,
		FuelPrice:      row.FuelPrice,
		DaysPerWeek:    int32(row.DaysPerWeek),
		AnnualCostRp:   row.AnnualCost,
		AnnualMinutes:  row.AnnualMinutes,
		AnnualHours:    float64(row.AnnualMinutes) / 60,
		AnnualWorkdays: float64(row.AnnualMinutes) / (60 * 8),
		CreatedAt:      row.CreatedAt.Format(time.RFC3339),
	}
	return commute, nil
}

func (s *CommuteService) DeleteCommute(ctx context.Context, id uuid.UUID) error {
	return s.store.DeleteCommute(ctx, id)
}

func coordsToLineString(coords [][]float64) orb.LineString {
	points := make([]orb.Point, len(coords))
	for i, c := range coords {
		points[i] = orb.Point{c[0], c[1]}
	}
	return orb.LineString(points)
}
