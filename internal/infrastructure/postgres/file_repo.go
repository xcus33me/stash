package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xcus33me/stash/internal/domain/file"
)

type fileRepo struct {
	db *pgxpool.Conn
}

func NewFileRepo(db *pgxpool.Conn) file.FileRepository {
	return &fileRepo{db}
}

func (r *fileRepo) Create(ctx context.Context, f *file.File) error {
	query := `
		INSERT INTO files
        	(id, room_id, original_name, uploader_alias, size_bytes, mime_type, storage_key, sha256, expires_at,
        VALUES
            ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.Exec(
		ctx, query, f.ID, f.RoomID, f.OriginalName,
		f.UploaderAlias, f.SizeBytes, f.MimeType,
		f.StorageKey, f.SHA256, f.ExpiresAt, f.UploadedAt,
	)

	return err
}

func (r *fileRepo) GetByID(ctx context.Context, fileID uuid.UUID) (*File, error) {
	f := &file.File{}

	query := `
		SELECT id, room_id, original_name, uploader_alias, size_bytes, mime_type, storage_key, sha256, expires_at, uploaded_at, deleted_at
        FROM files
    	WHERE id = $1
	`

}

func (r *fileRepo) ListByRoomID(ctx context.Context, roomID uuid.UUID) ([]*File, error) {

}

func (r *fileRepo) SoftDelete(ctx context.Context, fileID uuid.UUID) error {

}

func (r *fileRepo) ListExpired(ctx context.Context) ([]*File, error) {

}

func (r *fileRepo) TotalSizeByRoom(ctx context.Context, roomID uuid.UUID) (int64, error) {

}
