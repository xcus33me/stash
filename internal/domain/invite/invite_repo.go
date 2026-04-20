package invite

import (
	"context"

	"github.com/google/uuid"
)

type InviteLinkRepository interface {
	Create(ctx context.Context, il *InviteLink) error
	GetByID(ctx context.Context, id uuid.UUID) (*InviteLink, error)
	ListByRoomID(ctx context.Context, roomID uuid.UUID) ([]*InviteLink, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
