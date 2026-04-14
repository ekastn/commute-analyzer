package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// orsClient implements RoutingService using the OpenRouteService API.
type orsClient struct {
	apiKey string
	client *http.Client
}

func NewORSClient(apiKey string, client *http.Client) RoutingService {
	return &orsClient{
		apiKey: apiKey,
		client: client,
	}
}

func (c *orsClient) GetRoute(ctx context.Context, profile string, start, end Coord) (*RouteResult, error) {
	url := fmt.Sprintf(
		"https://api.openrouteservice.org/v2/directions/%s?api_key=%s&start=%f,%f&end=%f,%f",
		profile, c.apiKey, start.Lng, start.Lat, end.Lng, end.Lat,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ors error: %d", resp.StatusCode)
	}

	var orsResp struct {
		Features []struct {
			Properties struct {
				Summary struct {
					Distance float64 `json:"distance"` // meters
					Duration float64 `json:"duration"` // seconds
				} `json:"summary"`
			} `json:"properties"`
			Geometry struct {
				Coordinates [][]float64 `json:"coordinates"`
			} `json:"geometry"`
		} `json:"features"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&orsResp); err != nil {
		return nil, err
	}

	if len(orsResp.Features) == 0 {
		return nil, fmt.Errorf("ors error: no route found")
	}

	feature := orsResp.Features[0]
	summary := feature.Properties.Summary

	return &RouteResult{
		DistanceKm:  summary.Distance / 1000,
		DurationMin: summary.Duration / 60,
		Coordinates: feature.Geometry.Coordinates,
	}, nil
}
