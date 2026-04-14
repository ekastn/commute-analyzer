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
	routing     RoutingService
	userService *UserService
}

func NewCommuteService(store store.Querier, routing RoutingService, userService *UserService) *CommuteService {
	return &CommuteService{store: store, routing: routing, userService: userService}
}

func (s *CommuteService) CreateCommute(ctx context.Context, req dto.CreateCommuteRequest) (*dto.Commute, error) {
	userID, err := s.userService.GetOrCreateUser(ctx, req.DeviceID)
	if err != nil {
		return nil, err
	}

	// Map vehicle to routing profile
	profile := "driving-car" // Default profile

	switch req.Vehicle {
	case "car":
		profile = "driving-car"
	case "motorcycle":
		profile = "cycling-regular"
	// case "bicycle": // Example for future expansion
	// 	profile = "cycling-regular"
	}

	route, err := s.routing.GetRoute(ctx, profile, Coord{Lat: req.HomeLat, Lng: req.HomeLng}, Coord{Lat: req.OfficeLat, Lng: req.OfficeLng})
	if err != nil {
		return nil, fmt.Errorf("fetch route: %w", err)
	}

	distanceKm := route.DistanceKm
	durationMin := route.DurationMin

	lineString := coordsToLineString(route.Coordinates)

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
	// Fetch existing commute to get current values
	existing, err := s.store.GetCommute(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fetch existing: %w", err)
	}

	// Merge values: use request values if provided, otherwise fall back to existing
	name := existing.Name
	if req.Name != nil {
		name = req.Name
	}

	vehicle := existing.Vehicle
	if req.Vehicle != nil {
		vehicle = *req.Vehicle
	}

	fuelPrice := existing.FuelPrice
	if req.FuelPrice != nil {
		fuelPrice = int32(*req.FuelPrice)
	}

	daysPerWeek := existing.DaysPerWeek
	if req.DaysPerWeek != nil {
		daysPerWeek = int16(*req.DaysPerWeek)
	}

	// Parse existing coordinates
	existingHomeGeom, _ := wkb.Unmarshal(existing.HomePoint.([]byte))
	existingOfficeGeom, _ := wkb.Unmarshal(existing.OfficePoint.([]byte))
	existingHome := existingHomeGeom.(orb.Point)
	existingOffice := existingOfficeGeom.(orb.Point)

	// Determine if coordinates changed
	coordsChanged := false
	var homeLat, homeLng, officeLat, officeLng float64

	if req.HomeLat != nil && req.HomeLng != nil && req.OfficeLat != nil && req.OfficeLng != nil {
		homeLat, homeLng = *req.HomeLat, *req.HomeLng
		officeLat, officeLng = *req.OfficeLat, *req.OfficeLng
		if homeLat != existingHome.Lat() || homeLng != existingHome.Lon() ||
			officeLat != existingOffice.Lat() || officeLng != existingOffice.Lon() {
			coordsChanged = true
		}
	} else {
		// No coordinate update requested, use existing
		homeLat, homeLng = existingHome.Lat(), existingHome.Lon()
		officeLat, officeLng = existingOffice.Lat(), existingOffice.Lon()
	}

	// Route calculation: conditional ORS/RoutingService call
	distanceKm := existing.DistanceKm
	durationMin := existing.DurationMin
	var routeLineString orb.LineString

	if coordsChanged {
		// Map vehicle to routing profile
		profile := "driving-car"
		if vehicle == "motorcycle" {
			profile = "cycling-regular"
		}

		route, err := s.routing.GetRoute(ctx, profile, Coord{Lat: homeLat, Lng: homeLng}, Coord{Lat: officeLat, Lng: officeLng})
		if err != nil {
			return nil, fmt.Errorf("fetch route: %w", err)
		}

		distanceKm = route.DistanceKm
		durationMin = route.DurationMin
		routeLineString = coordsToLineString(route.Coordinates)
	} else {
		// Reuse existing route geometry
		existingRouteGeom, _ := wkb.Unmarshal(existing.RouteGeometry.([]byte))
		routeLineString = existingRouteGeom.(orb.LineString)
	}

	// Recalculate costs (always, since vehicle/fuel/days may have changed)
	efficiency := map[string]float64{"car": 10.0, "motorcycle": 2.5}[vehicle]
	roundTrip := distanceKm * 2
	dailyCost := (roundTrip * efficiency / 100) * float64(fuelPrice)
	annualCost := int64(dailyCost * float64(daysPerWeek) * 52.142857)
	annualMinutes := int64(durationMin*2) * int64(daysPerWeek) * 52

	// Marshal geometry
	homePoint := orb.Point{homeLng, homeLat}
	officePoint := orb.Point{officeLng, officeLat}
	homeWKB, _ := wkb.Marshal(homePoint)
	officeWKB, _ := wkb.Marshal(officePoint)
	routeWKB, _ := wkb.Marshal(routeLineString)

	// Update database
	row, err := s.store.UpdateCommute(ctx, store.UpdateCommuteParams{
		ID:            id,
		Name:          name,
		Vehicle:       &vehicle,
		FuelPrice:     &fuelPrice,
		DaysPerWeek:   &daysPerWeek,
		HomePoint:     homeWKB,
		OfficePoint:   officeWKB,
		RouteGeometry: routeWKB,
		DistanceKm:    &distanceKm,
		DurationMin:   &durationMin,
		AnnualCost:    &annualCost,
		AnnualMinutes: &annualMinutes,
	})
	if err != nil {
		return nil, fmt.Errorf("update: %w", err)
	}

	return &dto.Commute{
		ID:             existing.ID.String(),
		Name:           safeString(name),
		HomeLng:        homeLng,
		HomeLat:        homeLat,
		OfficeLng:      officeLng,
		OfficeLat:      officeLat,
		RouteGeometry:  &routeLineString,
		DistanceKm:     distanceKm,
		DurationMin:    durationMin,
		Vehicle:        vehicle,
		FuelPrice:      fuelPrice,
		DaysPerWeek:    int32(daysPerWeek),
		AnnualCostRp:   row.AnnualCost,
		AnnualMinutes:  row.AnnualMinutes,
		AnnualHours:    float64(row.AnnualMinutes) / 60,
		AnnualWorkdays: float64(row.AnnualMinutes) / (60 * 8),
		CreatedAt:      row.CreatedAt.Format(time.RFC3339),
	}, nil
}

func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
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
