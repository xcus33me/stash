package room

import (
	"time"

	"github.com/google/uuid"
)

type Options struct {
	MaxSizeBytes *int64
	FileTTL      *time.Duration
}

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

func NewRoom(title string, description *string, ownerTokenHash string, options Options) *Room {
	return &Room{
		ID:             uuid.New(),
		Title:          title,
		Description:    description,
		OwnerTokenHash: ownerTokenHash,
		Locked:         false,
		MaxSizeBytes:   options.MaxSizeBytes,
		FileTTL:        options.FileTTL,
		CreatedAt:      time.Now().UTC(),
	}
}

func (r *Room) Lock()   { r.Locked = true }
func (r *Room) Unlock() { r.Locked = false }
