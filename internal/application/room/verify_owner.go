package roomapp

import (
	"context"

	"github.com/google/uuid"
	"github.com/xcus33me/stash/internal/domain/room"
)

type VerifyOwnerInput struct {
	RoomID   uuid.UUID
	RawToken string
}

func (uc *Usecase) VerifyOwner(ctx context.Context, input VerifyOwnerInput) error {
	r, err := uc.loadRoom(ctx, input.RoomID)
	if err != nil {
		return err
	}

	if !uc.tokens.Verify(input.RawToken, r.OwnerTokenHash) {
		return room.ErrInvalidOwner
	}

	return nil
}
