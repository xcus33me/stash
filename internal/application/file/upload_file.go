package fileapp

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/xcus33me/stash/internal/domain/file"
	"github.com/xcus33me/stash/internal/domain/room"
)

const defaultMimeType = "application/octet-stream"

type UploadFileInput struct {
	RoomID        uuid.UUID
	OriginalName  string
	UploaderAlias *string
	MimeType      string
	SizeBytes     int64
	Reader        io.Reader
}

type UploadFileOutput struct {
	FileID    uuid.UUID
	SHA256    string
	ExpiresAt *time.Time
}

func (uc *Usecase) UploadFile(ctx context.Context, input UploadFileInput) (UploadFileOutput, error) {
	logger := uc.log.With("op", "upload_file", "room_id", input.RoomID)

	if input.OriginalName == "" {
		return UploadFileOutput{}, file.ErrNameRequired
	}
	if input.SizeBytes <= 0 {
		return UploadFileOutput{}, file.ErrInvalidSize
	}
	if input.Reader == nil {
		return UploadFileOutput{}, file.ErrInvalidSize
	}

	r, err := uc.loadRoom(ctx, input.RoomID)
	if err != nil {
		return UploadFileOutput{}, err
	}
	if r.Locked {
		return UploadFileOutput{}, room.ErrLocked
	}

	if r.MaxSizeBytes != nil {
		used, err := uc.files.TotalSizeByRoom(ctx, r.ID)
		if err != nil {
			return UploadFileOutput{}, fmt.Errorf("total size by room: %w", err)
		}
		if used+input.SizeBytes > *r.MaxSizeBytes {
			return UploadFileOutput{}, file.ErrQuotaExceeded
		}
	}

	mime := input.MimeType
	if mime == "" {
		mime = defaultMimeType
	}

	fileID := uuid.New()
	storageKey := fmt.Sprintf("%s/%s", r.ID, fileID)

	hasher := sha256.New()
	tee := io.TeeReader(input.Reader, hasher)

	if err := uc.storage.Upload(ctx, storageKey, tee, input.SizeBytes, mime); err != nil {
		return UploadFileOutput{}, fmt.Errorf("storage upload: %w", err)
	}

	sum := hex.EncodeToString(hasher.Sum(nil))

	var expiresAt *time.Time
	if r.FileTTL != nil {
		t := time.Now().UTC().Add(*r.FileTTL)
		expiresAt = &t
	}

	f := &file.File{
		ID:            fileID,
		RoomID:        r.ID,
		OriginalName:  input.OriginalName,
		UploaderAlias: input.UploaderAlias,
		SizeBytes:     input.SizeBytes,
		MimeType:      mime,
		StorageKey:    storageKey,
		SHA256:        sum,
		ExpiresAt:     expiresAt,
		UploadedAt:    time.Now().UTC(),
	}

	if err := uc.files.Create(ctx, f); err != nil {
		if delErr := uc.storage.Delete(ctx, storageKey); delErr != nil {
			logger.ErrorContext(ctx, "orphaned blob after failed db insert",
				"err", delErr,
				"storage_key", storageKey,
			)
		}
		return UploadFileOutput{}, fmt.Errorf("create file: %w", err)
	}

	logger.InfoContext(ctx, "file uploaded",
		"file_id", f.ID,
		"size_bytes", f.SizeBytes,
		"has_ttl", f.ExpiresAt != nil,
	)

	return UploadFileOutput{
		FileID:    f.ID,
		SHA256:    f.SHA256,
		ExpiresAt: f.ExpiresAt,
	}, nil
}
