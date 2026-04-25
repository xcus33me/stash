package fileapp

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/xcus33me/stash/internal/domain/room"
)

type DeleteFileInput struct {
	FileID     uuid.UUID
	OwnerToken string
}

func (uc *Usecase) DeleteFile(ctx context.Context, input DeleteFileInput) error {
	logger := uc.log.With("op", "delete_file", "file_id", input.FileID)

	f, err := uc.loadFile(ctx, input.FileID)
	if err != nil {
		return err
	}
	if f.IsDeleted() {
		return nil
	}

	r, err := uc.loadRoom(ctx, f.RoomID)
	if err != nil {
		return err
	}
	if !uc.tokens.Verify(input.OwnerToken, r.OwnerTokenHash) {
		return room.ErrInvalidOwner
	}

	if err := uc.storage.Delete(ctx, f.StorageKey); err != nil {
		logger.ErrorContext(ctx, "storage delete failed, continuing with soft delete",
			"err", err,
			"storage_key", f.StorageKey,
		)
	}

	if err := uc.files.SoftDelete(ctx, f.ID); err != nil {
		return fmt.Errorf("soft delete: %w", err)
	}

	logger.InfoContext(ctx, "file deleted", "room_id", f.RoomID)
	return nil
}
