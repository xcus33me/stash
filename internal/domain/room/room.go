package room

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID             uuid.UUID
	Title          string
	Description    *string
	OwnerTokenHash string
	Locked         bool
	MaxSizeBytes   *int64
	FileTTL        *time.Duration
	CreatedAt      time.Time
}

func (r *Room) Lock() { r.Locked = true }
