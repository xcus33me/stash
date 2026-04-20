package file

import (
	"context"

	"github.com/google/uuid"
)

type FileRepository interface {
	Create(ctx context.Context, file *File) error
	GetByID(ctx context.Context, fileID uuid.UUID) (*File, error)
	ListByRoomID(ctx context.Context, roomID uuid.UUID) ([]*File, error)
	SoftDelete(ctx context.Context, fileID uuid.UUID) error
	ListExpired(ctx context.Context) ([]*File, error)
	TotalSizeByRoom(ctx context.Context, roomID uuid.UUID) (int64, error)
}
