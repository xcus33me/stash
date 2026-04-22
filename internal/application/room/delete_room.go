package roomapp

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type DeleteRoomInput struct {
	RoomID     uuid.UUID
	OwnerToken string
}

func (uc *Usecase) DeleteRoom(ctx context.Context, input DeleteRoomInput) error {
	logger := uc.log.With("op", "delete_room")

	if err := uc.VerifyOwner(ctx, VerifyOwnerInput{input.RoomID, input.OwnerToken}); err != nil {
		return err
	}

	files, err := uc.files.ListByRoom(ctx, input.RoomID)
	if err != nil {
		return fmt.Errorf("list files: %w", err)
	}

	for _, f := range files {
		if err := uc.storage.Delete(ctx, f.StorageKey); err != nil {
			// Если ошибка при удалении файлов комнаты - пока просто логируем
			logger.ErrorContext(ctx, "error while deleting file attached to the room",
				"err", err,
				"room_id", f.RoomID,
				"file_id", f.ID,
			)
		}
	}

	if err := uc.files.DeleteByRoom(ctx, input.RoomID); err != nil {
		return fmt.Errorf("delete files: %w", err)
	}

	if err := uc.rooms.Delete(ctx, input.RoomID); err != nil {
		return fmt.Errorf("delete room: %w", err)
	}

	return nil
}
