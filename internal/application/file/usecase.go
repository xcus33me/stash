package fileapp

import (
	"context"
	"io"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/xcus33me/stash/internal/domain/file"
	"github.com/xcus33me/stash/internal/domain/room"
)

type TokenVerifier interface {
	Verify(raw, hash string) bool
}

type BlobStorage interface {
	Upload(ctx context.Context, key string, r io.Reader, size int64, mimeType string) error
	Delete(ctx context.Context, key string) error
	PresignedURL(ctx context.Context, key string, ttl time.Duration) (string, error)
}

type Usecase struct {
	files   file.FileRepository
	rooms   room.RoomRepository
	storage BlobStorage
	tokens  TokenVerifier

	downloadTTL time.Duration
	log         *slog.Logger
}

func NewUsecase(
	files file.FileRepository,
	rooms room.RoomRepository,
	storage BlobStorage,
	tokens TokenVerifier,
	downloadTTL time.Duration,
	logger *slog.Logger,
) *Usecase {
	return &Usecase{
		files:       files,
		rooms:       rooms,
		storage:     storage,
		tokens:      tokens,
		downloadTTL: downloadTTL,
		log:         logger.With("component", "fileapp"),
	}
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

func (uc *Usecase) loadFile(ctx context.Context, id uuid.UUID) (*file.File, error) {
	f, err := uc.files.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if f == nil {
		return nil, file.ErrNotFound
	}
	return f, nil
}
