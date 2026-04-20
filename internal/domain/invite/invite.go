package invite

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleReadOnly   Role = "read_only"
	RoleUploadOnly Role = "upload_only"
	RoleReadWrite  Role = "read_write"
	RoleAdmin      Role = "admin"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleReadOnly, RoleUploadOnly, RoleReadWrite, RoleAdmin:
		return true
	}
	return false
}

type InviteLink struct {
	ID        uuid.UUID
	RoomID    uuid.UUID
	Role      Role
	ExpiresAt *time.Time
	CreatedAt time.Time
}

func (i *InviteLink) IsExpired() bool {
	return i.ExpiresAt != nil && time.Now().After(*i.ExpiresAt)
}
