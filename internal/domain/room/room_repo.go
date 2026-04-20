package room

import (
	"context"

	"github.com/google/uuid"
)

type RoomRepository interface {
	Create(ctx context.Context, room *Room) error
	GetByID(ctx context.Context, id uuid.UUID) (*Room, error)
	Update(ctx context.Context, room *Room) error
	Delete(ctx context.Context, id uuid.UUID) error
}
