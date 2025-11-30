package service

import (
	"context"
)

type ORSClient interface {
	FetchRoute(ctx context.Context, profile string, homeLng, homeLat, officeLng, officeLat float64) (*ORSResponse, error)
}
