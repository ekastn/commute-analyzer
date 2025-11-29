package service

import (
	"context"
	"errors"

	"github.com/ekastn/commute-analyzer/internal/store"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserService struct {
	store store.Querier
}

func NewUserService(store store.Querier) *UserService {
	return &UserService{store: store}
}

func (s *UserService) GetOrCreateUser(ctx context.Context, deviceID string) (uuid.UUID, error) {
	id, err := s.store.GetUserByDeviceId(ctx, deviceID)
	if err == nil {
		return id, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return s.store.CreateUser(ctx, deviceID)
	}

	return uuid.Nil, err
}
