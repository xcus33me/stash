package roomapp

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/xcus33me/stash/internal/domain/file"
	"github.com/xcus33me/stash/internal/domain/room"
)

type TokenGenerator interface {
	Generate() (raw string, hash string, err error)
	Verify(raw, hash string) bool
}

type BlobStorage interface {
	Delete(ctx context.Context, key string) error
}

type fileLister interface {
	ListByRoom(ctx context.Context, roomID uuid.UUID) ([]*file.File, error)
	DeleteByRoom(ctx context.Context, roomID uuid.UUID) error
}

type Usecase struct {
	rooms   room.RoomRepository
	files   fileLister
	storage BlobStorage
	tokens  TokenGenerator

	log *slog.Logger
}

func NewUsecase(
	rooms room.RoomRepository,
	files fileLister,
	storage BlobStorage,
	tokens TokenGenerator,
	logger *slog.Logger,
) *Usecase {
	return &Usecase{rooms, files, storage, tokens, logger.With("component", "roomapp")}
}
