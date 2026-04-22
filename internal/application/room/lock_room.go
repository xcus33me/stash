package roomapp

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type LockUnlockRoomInput struct {
	RoomID     uuid.UUID
	OwnerToken string
}

func (uc *Usecase) LockRoom(ctx context.Context, input LockUnlockRoomInput) error {
	logger := uc.log.With("op", "lock_room")

	if err := uc.VerifyOwner(ctx, VerifyOwnerInput{input.RoomID, input.OwnerToken}); err != nil {
		return err
	}

	r, err := uc.loadRoom(ctx, input.RoomID)
	if err != nil {
		return err
	}

	r.Lock()

	if err := uc.rooms.Update(ctx, r); err != nil {
		return fmt.Errorf("failed to update room: %w", err)
	}

	logger.InfoContext(ctx, "room successfully locked", "room_id", r.ID)

	return nil
}

func (uc *Usecase) UnlockRoom(ctx context.Context, input LockUnlockRoomInput) error {

	logger := uc.log.With("op", "unlock_room")

	if err := uc.VerifyOwner(ctx, VerifyOwnerInput{input.RoomID, input.OwnerToken}); err != nil {
		return err
	}

	r, err := uc.loadRoom(ctx, input.RoomID)
	if err != nil {
		return err
	}

	r.Unlock()

	if err := uc.rooms.Update(ctx, r); err != nil {
		return fmt.Errorf("failed to update room: %w", err)
	}

	logger.InfoContext(ctx, "room successfully unlocked", "room_id", r.ID)

	return nil
}
