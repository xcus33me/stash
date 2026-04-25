package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xcus33me/stash/internal/domain/invite"
)

type inviteRepo struct {
	db *pgxpool.Pool
}

func NewInviteRepo(db *pgxpool.Pool) invite.InviteLinkRepository {
	return &inviteRepo{db}
}

func (r *inviteRepo) Create(ctx context.Context, il *invite.InviteLink) error {
	query := `
		INSERT INTO invite_links
			(id, room_id, role, expires_at, created_at)
		VALUES
			($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(
		ctx, query,
		il.ID, il.RoomID, il.Role, il.ExpiresAt, il.CreatedAt,
	)

	return err
}

func (r *inviteRepo) GetByID(ctx context.Context, id uuid.UUID) (*invite.InviteLink, error) {
	il := &invite.InviteLink{}

	query := `
		SELECT id, room_id, role, expires_at, created_at
		FROM invite_links
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&il.ID, &il.RoomID, &il.Role, &il.ExpiresAt, &il.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return il, nil
}

func (r *inviteRepo) ListByRoomID(ctx context.Context, roomID uuid.UUID) ([]*invite.InviteLink, error) {
	query := `
		SELECT id, room_id, role, expires_at, created_at
		FROM invite_links
		WHERE room_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []*invite.InviteLink
	for rows.Next() {
		il := &invite.InviteLink{}
		if err := rows.Scan(
			&il.ID, &il.RoomID, &il.Role, &il.ExpiresAt, &il.CreatedAt,
		); err != nil {
			return nil, err
		}
		links = append(links, il)
	}

	return links, rows.Err()
}

func (r *inviteRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM invite_links WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	return err
}
