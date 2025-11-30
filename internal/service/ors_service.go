package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ORSResponse struct {
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

type orsClient struct {
	apiKey string
	client *http.Client
}

func NewORSClient(apiKey string, client *http.Client) ORSClient {
	return &orsClient{
		apiKey: apiKey,
		client: client,
	}
}

func (c *orsClient) FetchRoute(ctx context.Context, profile string, homeLng, homeLat, officeLng, officeLat float64) (*ORSResponse, error) {
	url := fmt.Sprintf(
		"https://api.openrouteservice.org/v2/directions/%s?api_key=%s&start=%f,%f&end=%f,%f",
		profile, c.apiKey, homeLng, homeLat, officeLng, officeLat,
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

	var data ORSResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}
