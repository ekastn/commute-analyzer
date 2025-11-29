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

	route, err := s.ors.FetchRoute(ctx, req.HomeLng, req.HomeLat, req.OfficeLng, req.OfficeLat)
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
	if name == "" {
		name = fmt.Sprintf("Rute %s (%.1f km)", req.Vehicle, distanceKm)
	}

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
		FuelPrice:     req.FuelPrice,
		DaysPerWeek:   req.DaysPerWeek,
		AnnualCost:    annualCost,
		AnnualMinutes: annualMinutes,
	})
	if err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}

	return &dto.Commute{
		ID:             row.ID,
		Name:           name,
		HomeLng:        req.HomeLng,
		HomeLat:        req.HomeLat,
		OfficeLng:      req.OfficeLng,
		OfficeLat:      req.OfficeLat,
		DistanceKm:     distanceKm,
		DurationMin:    durationMin,
		Vehicle:        req.Vehicle,
		FuelPrice:      req.FuelPrice,
		DaysPerWeek:    req.DaysPerWeek,
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

		commutes[i] = dto.Commute{
			ID:             r.ID,
			Name:           name,
			HomeLng:        homePoint.Lon(),
			HomeLat:        homePoint.Lat(),
			OfficeLng:      officePoint.Lon(),
			OfficeLat:      officePoint.Lat(),
			DistanceKm:     r.DistanceKm,
			DurationMin:    r.DurationMin,
			Vehicle:        r.Vehicle,
			FuelPrice:      r.FuelPrice,
			DaysPerWeek:    r.DaysPerWeek,
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

	commute := &dto.Commute{
		ID:             row.ID,
		Name:           *row.Name,
		HomeLng:        hPoint.Lon(),
		HomeLat:        hPoint.Lat(),
		OfficeLng:      oPoint.Lon(),
		OfficeLat:      oPoint.Lat(),
		DistanceKm:     row.DistanceKm,
		DurationMin:    row.DurationMin,
		Vehicle:        row.Vehicle,
		FuelPrice:      row.FuelPrice,
		DaysPerWeek:    row.DaysPerWeek,
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
