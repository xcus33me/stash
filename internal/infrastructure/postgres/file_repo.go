package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xcus33me/stash/internal/domain/file"
)

type fileRepo struct {
	db *pgxpool.Pool
}

func NewFileRepo(db *pgxpool.Pool) file.FileRepository {
	return &fileRepo{db}
}

func (r *fileRepo) Create(ctx context.Context, f *file.File) error {
	query := `
		INSERT INTO files
			(id, room_id, original_name, uploader_alias, size_bytes, mime_type, storage_key, sha256, expires_at, uploaded_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.Exec(
		ctx, query,
		f.ID, f.RoomID, f.OriginalName,
		f.UploaderAlias, f.SizeBytes, f.MimeType,
		f.StorageKey, f.SHA256, f.ExpiresAt, f.UploadedAt,
	)

	return err
}

func (r *fileRepo) GetByID(ctx context.Context, fileID uuid.UUID) (*file.File, error) {
	f := &file.File{}

	query := `
		SELECT id, room_id, original_name, uploader_alias, size_bytes, mime_type, storage_key, sha256, expires_at, uploaded_at, deleted_at
		FROM files
		WHERE id = $1 AND deleted_at IS NULL
	`

	err := r.db.QueryRow(ctx, query, fileID).Scan(
		&f.ID, &f.RoomID, &f.OriginalName,
		&f.UploaderAlias, &f.SizeBytes, &f.MimeType,
		&f.StorageKey, &f.SHA256, &f.ExpiresAt,
		&f.UploadedAt, &f.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (r *fileRepo) ListByRoomID(ctx context.Context, roomID uuid.UUID) ([]*file.File, error) {
	query := `
		SELECT id, room_id, original_name, uploader_alias, size_bytes, mime_type, storage_key, sha256, expires_at, uploaded_at, deleted_at
		FROM files
		WHERE room_id = $1 AND deleted_at IS NULL
		ORDER BY uploaded_at DESC
	`

	rows, err := r.db.Query(ctx, query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*file.File
	for rows.Next() {
		f := &file.File{}
		if err := rows.Scan(
			&f.ID, &f.RoomID, &f.OriginalName,
			&f.UploaderAlias, &f.SizeBytes, &f.MimeType,
			&f.StorageKey, &f.SHA256, &f.ExpiresAt,
			&f.UploadedAt, &f.DeletedAt,
		); err != nil {
			return nil, err
		}
		files = append(files, f)
	}

	return files, rows.Err()
}

func (r *fileRepo) SoftDelete(ctx context.Context, fileID uuid.UUID) error {
	query := `
		UPDATE files
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Exec(ctx, query, time.Now(), fileID)
	return err
}

func (r *fileRepo) ListExpired(ctx context.Context) ([]*file.File, error) {
	query := `
		SELECT id, room_id, original_name, uploader_alias, size_bytes, mime_type, storage_key, sha256, expires_at, uploaded_at, deleted_at
		FROM files
		WHERE expires_at IS NOT NULL AND expires_at < $1 AND deleted_at IS NULL
	`

	rows, err := r.db.Query(ctx, query, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*file.File
	for rows.Next() {
		f := &file.File{}
		if err := rows.Scan(
			&f.ID, &f.RoomID, &f.OriginalName,
			&f.UploaderAlias, &f.SizeBytes, &f.MimeType,
			&f.StorageKey, &f.SHA256, &f.ExpiresAt,
			&f.UploadedAt, &f.DeletedAt,
		); err != nil {
			return nil, err
		}
		files = append(files, f)
	}

	return files, rows.Err()
}

func (r *fileRepo) TotalSizeByRoom(ctx context.Context, roomID uuid.UUID) (int64, error) {
	query := `
		SELECT COALESCE(SUM(size_bytes), 0)
		FROM files
		WHERE room_id = $1 AND deleted_at IS NULL
	`

	var total int64
	err := r.db.QueryRow(ctx, query, roomID).Scan(&total)
	return total, err
}
