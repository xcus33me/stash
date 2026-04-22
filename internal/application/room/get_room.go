package roomapp

import (
	"context"

	"github.com/google/uuid"
	"github.com/xcus33me/stash/internal/domain/room"
)

type GetRoomInput struct {
	RoomID uuid.UUID
}

type GetRoomOutput struct {
	ID          uuid.UUID
	Title       string
	Description *string
	Locked      bool
}

func (uc *Usecase) GetRoom(ctx context.Context, input GetRoomInput) (GetRoomOutput, error) {
	logger := uc.log.With("op", "get_room")

	r, err := uc.loadRoom(ctx, input.RoomID)
	if err != nil {
		return GetRoomOutput{}, err
	}

	logger.InfoContext(ctx, "succesfully got room", "room_id", r.ID)

	return GetRoomOutput{
		ID:          r.ID,
		Title:       r.Title,
		Description: r.Description,
		Locked:      r.Locked,
	}, nil
}

func (uc *Usecase) loadRoom(ctx context.Context, id uuid.UUID) (*room.Room, error) {
	r, err := uc.rooms.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, room.ErrNotFound
	}
	return r, nil
}
