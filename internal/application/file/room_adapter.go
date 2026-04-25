package fileapp

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/xcus33me/stash/internal/domain/file"
)

// ListByRoom / DeleteByRoom satisfy roomapp.fileLister so that roomapp can
// cascade through fileapp without knowing about the repository directly.

func (uc *Usecase) ListByRoom(ctx context.Context, roomID uuid.UUID) ([]*file.File, error) {
	return uc.files.ListByRoomID(ctx, roomID)
}

func (uc *Usecase) DeleteByRoom(ctx context.Context, roomID uuid.UUID) error {
	items, err := uc.files.ListByRoomID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("list files: %w", err)
	}
	for _, f := range items {
		if err := uc.files.SoftDelete(ctx, f.ID); err != nil {
			return fmt.Errorf("soft delete %s: %w", f.ID, err)
		}
	}
	return nil
}
