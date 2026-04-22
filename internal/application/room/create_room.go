package roomapp

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/xcus33me/stash/internal/domain/room"
)

type CreateRoomInput struct {
	Title        string
	Description  *string
	MaxSizeBytes *int64
	FileTTL      *time.Duration
}

type CreateRoomOutput struct {
	RoomID     uuid.UUID
	OwnerToken string
}

func (uc *Usecase) CreateRoom(ctx context.Context, input CreateRoomInput) (CreateRoomOutput, error) {
	logger := uc.log.With("op", "create_room")

	if input.Title == "" {
		return CreateRoomOutput{}, room.ErrTitleRequired
	}

	raw, hash, err := uc.tokens.Generate()
	if err != nil {
		return CreateRoomOutput{}, fmt.Errorf("generate token error: %w", err)
	}

	r := room.NewRoom(input.Title, input.Description, hash, room.Options{
		MaxSizeBytes: input.MaxSizeBytes,
		FileTTL:      input.FileTTL,
	})

	if err := uc.rooms.Create(ctx, r); err != nil {
		return CreateRoomOutput{}, fmt.Errorf("create room error: %w", err)
	}

	logger.InfoContext(ctx, "room successfully created",
		"room_id", r.ID,
		"has_ttl", r.FileTTL != nil,
		"has_size_limit", r.MaxSizeBytes != nil,
	)

	return CreateRoomOutput{RoomID: r.ID, OwnerToken: raw}, nil
}
