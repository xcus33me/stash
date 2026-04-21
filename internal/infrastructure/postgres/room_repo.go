package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xcus33me/stash/internal/domain/room"
)

type roomRepo struct {
	db *pgxpool.Conn
}

func NewRoomRepo(db *pgxpool.Conn) room.RoomRepository {
	return &roomRepo{db}
}

func (r *roomRepo) Create(ctx context.Context, rm *room.Room) error {
	query := `
		INSERT INTO rooms
			(id, title, description, owner_token_hash, locked, max_size_bytes, file_ttl, created_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(
		ctx, query,
		rm.ID, rm.Title, rm.Description, rm.OwnerTokenHash,
		rm.Locked, rm.MaxSizeBytes, rm.FileTTL, rm.CreatedAt,
	)

	return err
}

func (r *roomRepo) GetByID(ctx context.Context, id uuid.UUID) (*room.Room, error) {
	rm := &room.Room{}

	query := `
		SELECT id, title, description, owner_token_hash, locked, max_size_bytes, file_ttl, created_at
		FROM rooms
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&rm.ID, &rm.Title, &rm.Description, &rm.OwnerTokenHash,
		&rm.Locked, &rm.MaxSizeBytes, &rm.FileTTL, &rm.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return rm, nil
}

func (r *roomRepo) Update(ctx context.Context, rm *room.Room) error {
	query := `
		UPDATE rooms
		SET title = $1, description = $2, locked = $3, max_size_bytes = $4, file_ttl = $5
		WHERE id = $6
	`

	_, err := r.db.Exec(
		ctx, query,
		rm.Title, rm.Description, rm.Locked, rm.MaxSizeBytes, rm.FileTTL, rm.ID,
	)

	return err
}

func (r *roomRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM rooms WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	return err
}
