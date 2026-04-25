package fileapp

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/xcus33me/stash/internal/domain/file"
)

type DownloadFileInput struct {
	FileID uuid.UUID
}

type DownloadFileOutput struct {
	URL          string
	OriginalName string
	MimeType     string
	SizeBytes    int64
	ExpiresAt    *time.Time
}

func (uc *Usecase) DownloadFile(ctx context.Context, input DownloadFileInput) (DownloadFileOutput, error) {
	logger := uc.log.With("op", "download_file", "file_id", input.FileID)

	f, err := uc.loadFile(ctx, input.FileID)
	if err != nil {
		return DownloadFileOutput{}, err
	}
	if f.IsDeleted() {
		return DownloadFileOutput{}, file.ErrDeleted
	}
	if f.IsExpired() {
		return DownloadFileOutput{}, file.ErrExpired
	}

	url, err := uc.storage.PresignedURL(ctx, f.StorageKey, uc.downloadTTL)
	if err != nil {
		return DownloadFileOutput{}, fmt.Errorf("presigned url: %w", err)
	}

	logger.InfoContext(ctx, "presigned download url issued", "room_id", f.RoomID)

	return DownloadFileOutput{
		URL:          url,
		OriginalName: f.OriginalName,
		MimeType:     f.MimeType,
		SizeBytes:    f.SizeBytes,
		ExpiresAt:    f.ExpiresAt,
	}, nil
}
