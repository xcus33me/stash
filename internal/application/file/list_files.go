package fileapp

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ListFilesInput struct {
	RoomID uuid.UUID
}

type FileView struct {
	ID            uuid.UUID
	OriginalName  string
	UploaderAlias *string
	SizeBytes     int64
	MimeType      string
	UploadedAt    time.Time
	ExpiresAt     *time.Time
}

type ListFilesOutput struct {
	Files []FileView
}

func (uc *Usecase) ListFiles(ctx context.Context, input ListFilesInput) (ListFilesOutput, error) {
	if _, err := uc.loadRoom(ctx, input.RoomID); err != nil {
		return ListFilesOutput{}, err
	}

	items, err := uc.files.ListByRoomID(ctx, input.RoomID)
	if err != nil {
		return ListFilesOutput{}, fmt.Errorf("list files: %w", err)
	}

	views := make([]FileView, 0, len(items))
	for _, f := range items {
		if f.IsExpired() {
			continue
		}
		views = append(views, FileView{
			ID:            f.ID,
			OriginalName:  f.OriginalName,
			UploaderAlias: f.UploaderAlias,
			SizeBytes:     f.SizeBytes,
			MimeType:      f.MimeType,
			UploadedAt:    f.UploadedAt,
			ExpiresAt:     f.ExpiresAt,
		})
	}

	return ListFilesOutput{Files: views}, nil
}
